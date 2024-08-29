package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"wallet/entity"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
)

// ============== HELPER FUNCTIONS ==============
func findUserByID(id int) (entity.User, error) {
	for _, user := range entity.Users {
		if user.ID == id {
			return user, nil
		}
	}
	return entity.User{}, entity.ErrUserNotFoundInProducer
}

// ============== KAFKA RELATED FUNCTIONS ==============
func sendKafkaMessage(producer sarama.SyncProducer, ctx *gin.Context, fromID, toID int) error {
	// Retrieve 'nominal' as a string from the POST form data
	nominalStr := ctx.PostForm("nominal")

	// Convert 'nominal' from string to integer
	nominal, err := strconv.Atoi(nominalStr)
	if err != nil {
		return fmt.Errorf("failed to convert nominal to int: %w", err)
	}

	fromUser, err := findUserByID(fromID)
	if err != nil {
		return err
	}

	toUser, err := findUserByID(toID)
	if err != nil {
		return err
	}

	transaction := entity.Transaction{
		FromUserID: fromUser,
		ToUserID:   toUser,
		Nominal:    nominal,
		Timestamp:  time.Now(),
	}

	transactionJSON, err := json.Marshal(transaction)
	if err != nil {
		return fmt.Errorf("failed to marshal transaction: %w", err)
	}

	msg := &sarama.ProducerMessage{
		Topic: entity.Topic,
		Key:   sarama.StringEncoder(fromUser.ID),
		Value: sarama.StringEncoder(transactionJSON),
	}

	_, _, err = producer.SendMessage(msg)
	return err
}

func sendMessageHandler(producer sarama.SyncProducer) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		fromIDStr := ctx.PostForm("fromID")
		if fromIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "fromID is required"})
			return
		}
		fromID, err := strconv.Atoi(fromIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid fromID"})
			return
		}

		toIDStr := ctx.PostForm("toID")
		if toIDStr == "" {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "toID is required"})
			return
		}
		toID, err := strconv.Atoi(toIDStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"message": "Invalid toID"})
			return
		}

		err = sendKafkaMessage(producer, ctx, fromID, toID)
		if errors.Is(err, entity.ErrUserNotFoundInProducer) {
			ctx.JSON(http.StatusNotFound, gin.H{"message": "User not found"})
			return
		}
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{
				"message": err.Error(),
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{
			"message": "Notification sent successfully!",
		})
	}
}

func setupProducer() (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer([]string{entity.BrokerAddress},
		config)
	if err != nil {
		return nil, fmt.Errorf("failed to setup producer: %w", err)
	}
	return producer, nil
}

func main() {
	producer, err := setupProducer()
	if err != nil {
		log.Fatalf("failed to initialize producer: %v", err)
	}
	defer producer.Close()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.POST("/send", sendMessageHandler(producer))

	fmt.Printf("Kafka PRODUCER started at http://localhost%s\n",
		entity.ProducerPort)

	if err := router.Run(entity.ProducerPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
