package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"solution1/session-21/cmd/consumer"
	"solution1/session-21/entity"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupConsumerGroupWrite(ctx context.Context, store *consumer.NotificationStore, db *gorm.DB) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{entity.BrokerAddress}, entity.GroupWriteID, config)
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumerInstance := &consumer.Consumer{
		Store: store,
		DB:    db, // Pass the database connection to the consumer
	}

	for {
		err = consumerGroup.Consume(ctx, []string{entity.Topic}, consumerInstance)
		if err != nil {
			log.Printf("error from consumer: %v", err)
		}
		if ctx.Err() != nil {
			return
		}
	}
}

func saveNotificationToDB(db *gorm.DB, fromUserID int, message string) error {
	query := `INSERT INTO kafka_practice.notifications (user_id, message, timestamp) VALUES (?, ?, ?)`
	return db.Exec(query, fromUserID, message, time.Now()).Error
}

func main() {
	dsn := "postgresql://postgres:postgres@localhost:5433/postgres"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   "kafka_practice.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	store := &consumer.NotificationStore{
		Data: make(consumer.UserNotifications),
	}

	var notification entity.Notification

	// Extract the From user ID and message
	toUserID := notification.To.ID
	message := notification.Message

	// Save to the database
	if err := saveNotificationToDB(db, toUserID, message); err != nil {
		log.Printf("failed to save notification to database: %v", err)
	} else {
		log.Printf("Notification from user %s saved to the database", toUserID)
	}

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroupWrite(ctx, store, db)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		consumer.HandleNotifications(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", entity.GroupWriteID, entity.ConsumerWritePort)

	if err := router.Run(entity.ConsumerWritePort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
