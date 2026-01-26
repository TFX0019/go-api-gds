package auth

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name              string `gorm:"not null"`
	Email             string `gorm:"uniqueIndex;not null"`
	Password          string `gorm:"not null"`
	IsVerified        bool   `gorm:"default:false"`
	VerificationToken string
	ResetCode         string
	ResetCodeExpiry   time.Time
}
