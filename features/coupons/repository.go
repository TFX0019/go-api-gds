package coupons

import (
	"gorm.io/gorm"
)

type Repository interface {
	Get() (*Coupon, error)
	Update(coupon *Coupon) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Get() (*Coupon, error) {
	var coupon Coupon
	err := r.db.First(&coupon).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a default one if not exists
			defaultCoupon := Coupon{Code: "CUPON", Active: false}
			if createErr := r.db.Create(&defaultCoupon).Error; createErr != nil {
				return nil, createErr
			}
			return &defaultCoupon, nil
		}
		return nil, err
	}
	return &coupon, nil
}

func (r *repository) Update(coupon *Coupon) error {
	return r.db.Save(coupon).Error
}
