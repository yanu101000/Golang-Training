package entity

import (
	"time"
)

type Transaction struct {
	TransactionID   int       `gorm:"primaryKey"`
	UserID          int       `gorm:"not null"`
	Amount          float32   `gorm:"type:numeric(10,2);not null"`
	TransactionType string    `gorm:"not null"`
	TransactionDate time.Time `gorm:"not null"`
	Description     string
}
