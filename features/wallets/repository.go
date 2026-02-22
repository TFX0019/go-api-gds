package wallets

import (
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	AddCredits(userID uint, amount int, transactionType TransactionType, referenceID *string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) AddCredits(userID uint, amount int, transactionType TransactionType, referenceID *string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Log the credit transaction
		creditTx := &CreditTransaction{
			UserID:      userID,
			Amount:      amount,
			Type:        transactionType,
			ReferenceID: referenceID,
		}
		if err := tx.Create(creditTx).Error; err != nil {
			return err
		}

		// Update or Create the Wallet
		var wallet Wallet
		if err := tx.Where("user_id = ?", userID).First(&wallet).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				// Create wallet if it doesn't exist
				wallet = Wallet{
					UserID:       userID,
					Balance:      amount,
					LastRefillAt: time.Now(),
				}
				return tx.Create(&wallet).Error
			}
			return err
		}

		// Wallet exists, increment balance and update LastRefillAt
		wallet.Balance += amount
		wallet.LastRefillAt = time.Now()
		return tx.Save(&wallet).Error
	})
}
