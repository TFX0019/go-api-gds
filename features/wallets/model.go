package wallets

import (
	"time"

	"gorm.io/gorm"
)

type Wallet struct {
	ID           uint           `gorm:"primarykey" json:"id"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
	UserID       uint           `gorm:"not null;uniqueIndex"`
	Balance      int            `gorm:"default:0"`
	LastRefillAt time.Time      `gorm:"default:null"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
}

type TransactionType string

const (
	TransactionTypeSubscriptionRenewal TransactionType = "subscription_renewal"
	TransactionTypeSubscriptionUpgrade TransactionType = "subscription_upgrade"
	TransactionTypeImageGeneration     TransactionType = "image_generation"
)

// TODO: add payment gateway integration
type CreditTransaction struct {
	ID          uint            `gorm:"primarykey" json:"id"`
	UserID      uint            `gorm:"not null;index"`
	Amount      int             `gorm:"not null"`
	Type        TransactionType `gorm:"type:varchar(50);not null"`
	ReferenceID *string         `gorm:"type:varchar(255)"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}
