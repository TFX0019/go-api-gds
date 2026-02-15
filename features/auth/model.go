package auth

import (
	"time"

	"github.com/TFX0019/api-go-gds/features/subscriptions"
	"github.com/TFX0019/api-go-gds/features/wallets"
	"gorm.io/gorm"
)

type Role struct {
	gorm.Model
	Name        string `gorm:"unique;not null"`
	Description string
}

type Session struct {
	gorm.Model
	UserID    uint      `gorm:"not null"`
	Token     string    `gorm:"uniqueIndex;not null"` // Refresh token
	ExpiresAt time.Time `gorm:"not null"`
	IPAddress string
	UserAgent string
	IsValid   bool `gorm:"default:true"`
}

type User struct {
	gorm.Model
	Name              string `gorm:"not null"`
	Email             string `gorm:"uniqueIndex;not null"`
	Password          string `gorm:"not null"`
	IsVerified        bool   `gorm:"default:false"`
	IsActive          bool   `gorm:"default:true" json:"is_active"`
	VerificationToken string
	ResetCode         string
	ResetCodeExpiry   time.Time
	Avatar            *string
	Wallet            wallets.Wallet             `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Subscription      subscriptions.Subscription `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Roles             []Role                     `gorm:"many2many:user_roles;"`
}

type VerificationCode struct {
	gorm.Model
	Email     string    `gorm:"index;not null"`
	Code      string    `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
}
