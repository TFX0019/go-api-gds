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
