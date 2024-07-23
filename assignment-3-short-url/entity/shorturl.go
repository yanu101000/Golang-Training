package entity

import (
	"time"
)

type Url struct {
	ID          uint      `gorm:"primaryKey"`
	ShortUrl    string    `gorm:"size:255;not null"`
	OriginalUrl string    `gorm:"size:255;not null"`
	CreatedAt   time.Time `json:"created_at"`
}
