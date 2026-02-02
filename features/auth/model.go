package auth

import (
	"time"

	"github.com/TFX0019/api-go-gds/features/subscriptions"
	"github.com/TFX0019/api-go-gds/features/wallets"
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
	Avatar            *string
	Wallet            wallets.Wallet             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscription      subscriptions.Subscription `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
