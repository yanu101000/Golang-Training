package consumer

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"
	"time"

	"solution1/session-21/entity"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
type UserNotifications map[string][]entity.Notification

// NotificationStore holds user notifications and provides thread-safe access.
type NotificationStore struct {
	Data UserNotifications // Exported field to make it accessible
	mu   sync.RWMutex
}

// Add stores a new notification for a user.
func (ns *NotificationStore) Add(userID string, notification entity.Notification) {
	ns.mu.Lock()
	defer ns.mu.Unlock()
	ns.Data[userID] = append(ns.Data[userID], notification)
}

// Get retrieves all notifications for a specific user.
func (ns *NotificationStore) Get(userID string) []entity.Notification {
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
	for msg := range claim.Messages() {
		userID := string(msg.Key)
		var notification entity.Notification
		err := json.Unmarshal(msg.Value, &notification)
		if err != nil {
			log.Printf("failed to unmarshal notification: %v", err)
			continue
		}

		// Extract the From user ID and message
		fromUserID := notification.From.ID
		toUserID := notification.To.ID
		message := notification.Message

		// Save to the database
		if err := saveNotificationToDB(c.DB, fromUserID, toUserID, message); err != nil {
			log.Printf("failed to save notification to database: %v", err)
		} else {
			log.Printf("Notification from user %s saved to the database", fromUserID)
		}

		log.Printf("Consuming notification and adding it to storage: %v", notification)
		c.Store.Add(userID, notification)
		sess.MarkMessage(msg, "")
	}
	return nil
}

func saveNotificationToDB(db *gorm.DB, fromUserID int, toUserID int, message string) error {
	query := `INSERT INTO kafka_practice.notifications (fromUserID, toUserID, message, timestamp) VALUES (?, ?, ?, ?)`
	return db.Exec(query, fromUserID, toUserID, message, time.Now()).Error
}

func HandleNotifications(ctx *gin.Context, store *NotificationStore) {
	userID, err := getUserIDFromRequest(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"message": err.Error()})
		return
	}

	notes := store.Get(userID)
	if len(notes) == 0 {
		ctx.JSON(http.StatusOK,
			gin.H{
				"message":       "No notifications found for user",
				"notifications": []entity.Notification{},
			})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"notifications": notes})
}
