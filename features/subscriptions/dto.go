package subscriptions

type RevenueCatWebhook struct {
	Event struct {
		EventTimestampMs int64  `json:"event_timestamp_ms"`
		AppUserID        string `json:"app_user_id"`
		Type             string `json:"type"`
		ProductID        string `json:"product_id"`
		ExpirationAtMs   int64  `json:"expiration_at_ms,omitempty"`
		PurchasedAtMs    int64  `json:"purchased_at_ms,omitempty"`
		Environment      string `json:"environment"`
	} `json:"event"`
	APIVersion string `json:"api_version"`
}
