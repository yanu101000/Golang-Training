package main

import (
	"context"
	"fmt"
	"log"

	"solution1/session-21/cmd/consumer"
	"solution1/session-21/entity"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

func setupConsumerGroupRead(ctx context.Context, store *consumer.NotificationStore, db *gorm.DB) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{entity.BrokerAddress}, entity.GroupReadID, config)
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

func main() {
	dsn := "postgresql://postgres:postgres@localhost:5432/postgres"
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

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroupRead(ctx, store, db)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		consumer.HandleNotifications(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", entity.GroupReadID, entity.ConsumerReadPort)

	if err := router.Run(entity.ConsumerReadPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
