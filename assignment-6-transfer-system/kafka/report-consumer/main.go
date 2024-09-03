package consumer

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"wallet/entity"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// ============== HELPER FUNCTIONS ==============
var ErrNoMessagesFound = errors.New("no messages found")

func getUserIDFromRequest(ctx *gin.Context) (string, error) {
	userID := ctx.Param("userID")
	if userID == "" {
		return "", ErrNoMessagesFound
	}
	return userID, nil
}

// ====== NOTIFICATION STORAGE ======
type UserNotifications map[string][]entity.Transaction

// NotificationStore holds user notifications and provides thread-safe access.
type NotificationStore struct {
	Data UserNotifications // Exported field to make it accessible
	mu   sync.RWMutex
}

// Add stores a new transaction for a user.
func (ns *NotificationStore) Add(userID string, transaction entity.Transaction) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.Data[userID] = append(ns.Data[userID], transaction)
}

// Get retrieves all notifications for a specific user.
func (ns *NotificationStore) Get(userID string) []entity.Transaction {
	ns.mu.RLock()
	defer ns.mu.RUnlock()
	return ns.Data[userID]
}

// ============== KAFKA RELATED FUNCTIONS ==============
type Consumer struct {
	Store *NotificationStore
	DB    *gorm.DB
}

func (c *Consumer) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (c *Consumer) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Open the CSV file for writing; create it if it doesn't exist
	file, err := os.OpenFile("transaction-log.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("failed to open CSV file: %v", err)
		return err
	}
	defer file.Close()

	// Create a new CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	for msg := range claim.Messages() {
		userID := string(msg.Key)
		var transaction entity.Transaction
		err := json.Unmarshal(msg.Value, &transaction)
		if err != nil {
			log.Printf("failed to unmarshal transaction: %v", err)
			continue
		}

		// Extract the From user ID and nominal
		fromUserID := transaction.FromUserID.ID
		toUserID := transaction.ToUserID.ID
		nominal := transaction.Nominal

		// Save to the database
		// if err := saveNotificationToDB(c.DB, fromUserID, toUserID, nominal); err != nil {
		// 	log.Printf("failed to save transaction to database: %v", err)
		// } else {
		// 	log.Printf("Transaction from user %s saved to the database", fromUserID)
		// }

		// Log the transaction in the CSV file
		err = writeTransactionToCSV(writer, fromUserID, toUserID, nominal, time.Now())
		if err != nil {
			log.Printf("failed to write transaction to CSV: %v", err)
		}

		log.Printf("Consuming transaction and adding it to storage: %v", transaction)
		c.Store.Add(userID, transaction)
		sess.MarkMessage(msg, "")
	}
	return nil
}

// func saveNotificationToDB(db *gorm.DB, fromUserID int, toUserID int, nominal int) error {
// 	query := `INSERT INTO kafka_practice.transactions (fromUserID, toUserID, nominal, timestamp) VALUES (?, ?, ?, ?)`
// 	return db.Exec(query, fromUserID, toUserID, nominal, time.Now()).Error
// }

func writeTransactionToCSV(writer *csv.Writer, fromUserID int, toUserID int, nominal int, timestamp time.Time) error {
	// Convert data to strings and prepare the record for CSV writing
	record := []string{
		strconv.Itoa(fromUserID),
		strconv.Itoa(toUserID),
		strconv.Itoa(nominal),
		timestamp.Format(time.RFC3339), // Convert timestamp to string format
	}

	// Write the record to the CSV file
	if err := writer.Write(record); err != nil {
		return fmt.Errorf("failed to write record to CSV: %w", err)
	}

	return nil
}

func HandleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"nominal": err.Error()})
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		ctx.JSON(http.StatusOK,
			gin.H{
				"nominal":       "No notifications found for user",
				"notifications": []entity.Transaction{},
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}

func setupConsumerGroupRead(ctx context.Context, store *NotificationStore, db *gorm.DB) {
	config := sarama.NewConfig()

	consumerGroup, err := sarama.NewConsumerGroup([]string{entity.BrokerAddress}, entity.GroupReportID, config)
	if err != nil {
		log.Printf("initialization error: %v", err)
	}
	defer consumerGroup.Close()

	consumerInstance := &Consumer{
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
			TablePrefix:   "assignment-6.",
			SingularTable: false,
		},
	})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	store := &NotificationStore{
		Data: make(UserNotifications),
	}

	ctx, cancel := context.WithCancel(context.Background())
	go setupConsumerGroupRead(ctx, store, db)
	defer cancel()

	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.GET("/notifications/:userID", func(ctx *gin.Context) {
		HandleNotifications(ctx, store)
	})

	fmt.Printf("Kafka CONSUMER (Group: %s) 👥📥 "+
		"started at http://localhost%s\n", entity.GroupReportID, entity.ConsumerReportPort)

	if err := router.Run(entity.ConsumerReportPort); err != nil {
		log.Printf("failed to run the server: %v", err)
	}
}
