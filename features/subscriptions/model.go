package subscriptions

import (
	"time"
)

type SubscriptionStatus string

const (
	SubscriptionStatusActive      SubscriptionStatus = "active"
	SubscriptionStatusExpired     SubscriptionStatus = "expired"
	SubscriptionStatusGracePeriod SubscriptionStatus = "grace_period"
)

type Subscription struct {
	ID        uint               `gorm:"primarykey" json:"id"`
	UserID    uint               `gorm:"not null;uniqueIndex" json:"user_id"`
	ProductID string             `gorm:"type:varchar(255);not null" json:"product_id"`
	Status    SubscriptionStatus `gorm:"type:varchar(50);not null" json:"status"`
	ExpiresAt time.Time          `gorm:"not null" json:"expires_at"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
}

type Transaction struct {
	ID                    uint      `gorm:"primarykey" json:"id"`
	UserID                uint      `gorm:"not null" json:"user_id"`
	RevenueCatID          string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"revenuecat_id"`
	Type                  string    `gorm:"type:varchar(50);not null" json:"type"`
	ProductID             string    `gorm:"type:varchar(255);not null" json:"product_id"`
	Store                 string    `gorm:"type:varchar(50)" json:"store"`
	Environment           string    `gorm:"type:varchar(50)" json:"environment"`
	Currency              string    `gorm:"type:varchar(10)" json:"currency"`
	Price                 float64   `json:"price"`
	TransactionID         string    `gorm:"type:varchar(255)" json:"transaction_id"`
	OriginalTransactionID string    `gorm:"type:varchar(255)" json:"original_transaction_id"`
	EventTimestampMs      int64     `json:"event_timestamp_ms"`
	PurchasedAtMs         int64     `json:"purchased_at_ms"`
	ExpirationAtMs        int64     `json:"expiration_at_ms"`
	CreatedAt             time.Time `json:"created_at"`
	UpdatedAt             time.Time `json:"updated_at"`
}
