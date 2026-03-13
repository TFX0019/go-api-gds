package subscriptions

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetSubscriptionByUserID(userID uint) (*Subscription, error)
	UpsertSubscription(sub *Subscription) error
	CreateTransaction(t *Transaction) error
	GetAllTransactions(limit, offset int, search string) ([]TransactionResponse, int64, error)
	HavePurchasedPack80Credits(userID uint) (bool, error)
	GetActiveCoupon() (*string, error)
	GetUserEmail(userID uint) (*string, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db}
}

func (r *repository) GetSubscriptionByUserID(userID uint) (*Subscription, error) {
	var sub Subscription
	err := r.db.Where("user_id = ?", userID).First(&sub).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // Not found but it's not a severe error
		}
		return nil, err
	}
	return &sub, nil
}

func (r *repository) UpsertSubscription(sub *Subscription) error {
	var existing Subscription
	err := r.db.Where("user_id = ?", sub.UserID).First(&existing).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return err
	}

	if err == gorm.ErrRecordNotFound {
		return r.db.Create(sub).Error
	}

	sub.ID = existing.ID
	sub.CreatedAt = existing.CreatedAt
	return r.db.Save(sub).Error
}

func (r *repository) CreateTransaction(t *Transaction) error {
	return r.db.Create(t).Error
}

func (r *repository) GetAllTransactions(limit, offset int, search string) ([]TransactionResponse, int64, error) {
	var results []TransactionResponse
	var total int64

	query := r.db.Table("transactions").
		Select("transactions.*, users.name as user_name, users.email as user_email, plans.title as plan_name").
		Joins("LEFT JOIN users ON transactions.user_id = users.id").
		Joins("LEFT JOIN plans ON transactions.product_id = plans.product_id")

	if search != "" {
		searchPattern := "%" + search + "%"
		query = query.Where("users.name ILIKE ? OR users.email ILIKE ?", searchPattern, searchPattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("transactions.created_at DESC").
		Limit(limit).
		Offset(offset).
		Scan(&results).Error

	return results, total, err
}

func (r *repository) HavePurchasedPack80Credits(userID uint) (bool, error) {
	var count int64
	err := r.db.Model(&Transaction{}).Where("user_id = ? AND product_id = 'pack_80_credits' AND type IN ('INITIAL_PURCHASE', 'NON_RENEWING_PURCHASE')", userID).Count(&count).Error
	return count > 0, err
}

func (r *repository) GetActiveCoupon() (*string, error) {
	var code string
	err := r.db.Table("coupons").Where("active = ?", true).Select("code").First(&code).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil // No active coupon
		}
		return nil, err
	}
	return &code, nil
}

func (r *repository) GetUserEmail(userID uint) (*string, error) {
	var email string
	err := r.db.Table("users").Where("id = ?", userID).Select("email").First(&email).Error
	if err != nil {
		return nil, err
	}
	return &email, nil
}
