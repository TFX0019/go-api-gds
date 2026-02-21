package subscriptions

import (
	"log"
	"strconv"
	"time"
)

type Service interface {
	HandleRevenueCatWebhook(payload RevenueCatWebhook) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo}
}

func (s *service) HandleRevenueCatWebhook(payload RevenueCatWebhook) error {
	userID, err := strconv.ParseUint(payload.Event.AppUserID, 10, 32)
	if err != nil {
		log.Printf("[RevenueCat Webhook] Invalid AppUserID: %v", payload.Event.AppUserID)
		// Returning nil so RevenueCat doesn't keep retrying if it's an invalid format that we can't handle anyway.
		return nil
	}

	eventType := payload.Event.Type
	log.Printf("[RevenueCat Webhook] Received %s for user %d", eventType, userID)

	status := SubscriptionStatusActive
	switch eventType {
	case "EXPIRATION":
		status = SubscriptionStatusExpired
	case "BILLING_ISSUE":
		status = SubscriptionStatusGracePeriod
	case "CANCELLATION":
		// Cancellation means they turned off auto-renew. The plan is still "active" until ExpirationAtMs.
		// So we can still mark it as active, but wait if it's expired we will get EXPIRATION later.
		status = SubscriptionStatusActive
	}

	var expiresAt time.Time
	if payload.Event.ExpirationAtMs > 0 {
		expiresAt = time.UnixMilli(payload.Event.ExpirationAtMs)
	} else {
		expiresAt = time.Now().AddDate(1, 0, 0)
	}

	// We might have an existing subscription.
	// Make sure we only check valid plans, or maybe we just store the product ID sent by RevenueCat.
	sub := &Subscription{
		UserID:    uint(userID),
		ProductID: payload.Event.ProductID,
		Status:    status,
		ExpiresAt: expiresAt,
	}

	err = s.repo.UpsertSubscription(sub)
	if err != nil {
		log.Printf("[RevenueCat Webhook] Error upserting subscription for user %d: %v", userID, err)
		return err
	}

	return nil
}
