package entity

import (
	"time"
)

type User struct {
	ID          int       `gorm:"primaryKey" json:"id"`
	Name        string    `gorm:"type:varchar;not null" json:"name" binding:"required"`
	Email       string    `gorm:"type:varchar;uniqueIndex;not null" json:"email" binding:"required,email"`
	Mobilephone string    `gorm:"type:varchar;not null" json:"mobilephone" binding:"required"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
