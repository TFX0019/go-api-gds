package cronjobs

import (
	"log"
	"strconv"

	"github.com/TFX0019/api-go-gds/features/daily_credits"
	"github.com/TFX0019/api-go-gds/features/subscriptions"
	"github.com/TFX0019/api-go-gds/features/wallets"
	"github.com/TFX0019/api-go-gds/pkg/config"
	"gorm.io/gorm"
)

func CheckAndRefillCredits(db *gorm.DB) {
	// 1. Get REMOVE_CREDITS_FOR_GENERATION
	removeCreditsStr := config.GetEnv("REMOVE_CREDITS_FOR_GENERATION", "10")
	removeCredits, _ := strconv.Atoi(removeCreditsStr)
	if removeCredits == 0 {
		removeCredits = 10
	}

	// 2. Get Daily Credits Config
	var dc daily_credits.DailyCredit
	if err := db.First(&dc).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			dc = daily_credits.DailyCredit{Free: 3, Premium: 6}
			db.Create(&dc)
		} else {
			log.Printf("[Cronjob] Error getting daily credits: %v", err)
			return
		}
	}

	freeMax := dc.Free * removeCredits
	premiumMax := dc.Premium * removeCredits

	// 3. Batch process users
	var usersWithWallets []struct {
		UserID    uint
		Balance   int
		ProductID string
		Status    subscriptions.SubscriptionStatus
	}

	// We Join wallets and subscriptions to check user tier and balance
	err := db.Table("wallets").
		Select("wallets.user_id, wallets.balance, COALESCE(subscriptions.product_id, 'free_tier') as product_id, subscriptions.status").
		Joins("LEFT JOIN subscriptions ON wallets.user_id = subscriptions.user_id").
		Scan(&usersWithWallets).Error

	if err != nil {
		log.Printf("[Cronjob] Error getting wallets: %v", err)
		return
	}

	for _, uw := range usersWithWallets {
		isPremium := false
		if uw.ProductID != "free_tier" && uw.ProductID != "" && uw.Status == subscriptions.SubscriptionStatusActive {
			isPremium = true
		}

		newBalance := uw.Balance
		needsUpdate := false

		if isPremium {
			if uw.Balance < premiumMax {
				newBalance = premiumMax
				needsUpdate = true
			}
		} else {
			if uw.Balance < freeMax {
				newBalance = freeMax
				needsUpdate = true
			}
		}

		if needsUpdate {
			err := db.Model(&wallets.Wallet{}).Where("user_id = ?", uw.UserID).Update("balance", newBalance).Error
			if err != nil {
				log.Printf("[Cronjob] Error updating balance for user %d: %v", uw.UserID, err)
			} else {
				log.Printf("[Cronjob] Updated balance for user %d to %d", uw.UserID, newBalance)
			}
		}
	}
}
