package subscriptions

import (
	"log"
	"strconv"
	"time"

	"github.com/TFX0019/api-go-gds/features/wallets"
	"github.com/TFX0019/api-go-gds/pkg/config"
)

type Service interface {
	HandleRevenueCatWebhook(payload RevenueCatWebhook) error
	ListTransactions(page, limit int, search string) (*PaginatedTransactionResponse, error)
}

type service struct {
	repo        Repository
	walletsRepo wallets.Repository
}

func NewService(repo Repository, walletsRepo wallets.Repository) Service {
	return &service{repo, walletsRepo}
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

	// Save transaction log
	txn := &Transaction{
		UserID:                uint(userID),
		RevenueCatID:          payload.Event.ID,
		Type:                  eventType,
		ProductID:             payload.Event.ProductID,
		Store:                 payload.Event.Store,
		Environment:           payload.Event.Environment,
		Currency:              payload.Event.Currency,
		Price:                 payload.Event.Price,
		TransactionID:         payload.Event.TransactionID,
		OriginalTransactionID: payload.Event.OriginalTransactionID,
		EventTimestampMs:      payload.Event.EventTimestampMs,
		PurchasedAtMs:         payload.Event.PurchasedAtMs,
		ExpirationAtMs:        payload.Event.ExpirationAtMs,
	}

	err = s.repo.CreateTransaction(txn)
	if err != nil {
		log.Printf("[RevenueCat Webhook] Warning: could not log transaction for user %d: %v", userID, err)
		// We still return nil since subscription upsert succeeded, this is just a log.
	}

	if payload.Event.ProductID == "pack_80_credits" {
		if eventType == "NON_RENEWING_PURCHASE" || eventType == "INITIAL_PURCHASE" {
			removeCreditsStr := config.GetEnv("REMOVE_CREDITS_FOR_GENERATION", "10")
			removeCredits, _ := strconv.Atoi(removeCreditsStr)
			if removeCredits == 0 {
				removeCredits = 10
			}
			credits := 8 * removeCredits

			err = s.walletsRepo.AddCredits(uint(userID), credits, wallets.TransactionTypeAddCredits, &payload.Event.ID)
			if err != nil {
				log.Printf("[RevenueCat Webhook] Error adding credits for user %d: %v", userID, err)
				return err
			}
			log.Printf("[RevenueCat Webhook] Added %d credits to user %d from pack_80_credits", credits, userID)
		}
		return nil
	}

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
		expiresAt = time.UnixMilli(payload.Event.ExpirationAtMs).UTC()
	} else {
		expiresAt = time.Now().UTC().AddDate(1, 0, 0)
	}

	productID := payload.Event.ProductID
	if eventType == "EXPIRATION" {
		productID = "free_tier"
	}

	// We might have an existing subscription.
	// Make sure we only check valid plans, or maybe we just store the product ID sent by RevenueCat.
	sub := &Subscription{
		UserID:    uint(userID),
		ProductID: productID,
		Status:    status,
		ExpiresAt: expiresAt,
	}

	err = s.repo.UpsertSubscription(sub)
	if err != nil {
		log.Printf("[RevenueCat Webhook] Error upserting subscription for user %d: %v", userID, err)
		return err
	}

	// Add credits for purchase or renewal
	if eventType == "INITIAL_PURCHASE" || eventType == "RENEWAL" {
		creditsStr := config.GetEnv("ADD_CREDITS_SUBSCRIPTION", "100")
		credits, _ := strconv.Atoi(creditsStr)
		if credits == 0 {
			credits = 100 // Default fallback
		}

		// Use the correct transaction type
		var walletTxType wallets.TransactionType
		if eventType == "INITIAL_PURCHASE" {
			walletTxType = wallets.TransactionTypeSubscriptionUpgrade // or create a specific one
		} else {
			walletTxType = wallets.TransactionTypeSubscriptionRenewal
		}

		err = s.walletsRepo.AddCredits(uint(userID), credits, walletTxType, &payload.Event.ID)
		if err != nil {
			log.Printf("[RevenueCat Webhook] Error adding credits for user %d: %v", userID, err)
			return err
		}
		log.Printf("[RevenueCat Webhook] Added %d credits to user %d", credits, userID)
	}

	return nil
}

func (s *service) ListTransactions(page, limit int, search string) (*PaginatedTransactionResponse, error) {
	offset := (page - 1) * limit
	transactions, total, err := s.repo.GetAllTransactions(limit, offset, search)
	if err != nil {
		return nil, err
	}

	return &PaginatedTransactionResponse{
		Data:  transactions,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}
