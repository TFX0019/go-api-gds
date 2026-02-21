package subscriptions

import (
	"gorm.io/gorm"
)

type Repository interface {
	GetSubscriptionByUserID(userID uint) (*Subscription, error)
	UpsertSubscription(sub *Subscription) error
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
