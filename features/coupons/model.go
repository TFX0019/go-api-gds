package coupons

import "time"

type Coupon struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Code      string    `gorm:"type:varchar(255);not null;default:'CUPON'" json:"code"`
	Active    bool      `gorm:"default:false" json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Coupon) TableName() string {
	return "coupons"
}
