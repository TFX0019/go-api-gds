package subscriptions

import (
	"time"
)

type RevenueCatWebhook struct {
	Event struct {
		ID                    string  `json:"id"`
		EventTimestampMs      int64   `json:"event_timestamp_ms"`
		AppUserID             string  `json:"app_user_id"`
		Type                  string  `json:"type"`
		ProductID             string  `json:"product_id"`
		Store                 string  `json:"store"`
		Environment           string  `json:"environment"`
		Currency              string  `json:"currency"`
		Price                 float64 `json:"price"`
		TransactionID         string  `json:"transaction_id"`
		OriginalTransactionID string  `json:"original_transaction_id"`
		ExpirationAtMs        int64   `json:"expiration_at_ms,omitempty"`
		PurchasedAtMs         int64   `json:"purchased_at_ms,omitempty"`
	} `json:"event"`
	APIVersion string `json:"api_version"`
}

type TransactionResponse struct {
	ID                    uint      `json:"id"`
	UserID                uint      `json:"user_id"`
	UserName              string    `json:"user_name"`
	UserEmail             string    `json:"user_email"`
	RevenueCatID          string    `json:"revenuecat_id"`
	Type                  string    `json:"type"`
	ProductID             string    `json:"product_id"`
	Store                 string    `json:"store"`
	Environment           string    `json:"environment"`
	Currency              string    `json:"currency"`
	Price                 float64   `json:"price"`
	TransactionID         string    `json:"transaction_id"`
	OriginalTransactionID string    `json:"original_transaction_id"`
	EventTimestampMs      int64     `json:"event_timestamp_ms"`
	PurchasedAtMs         int64     `json:"purchased_at_ms"`
	ExpirationAtMs        int64     `json:"expiration_at_ms"`
	CreatedAt             time.Time `json:"created_at"`
}
