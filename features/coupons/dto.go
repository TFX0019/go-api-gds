package coupons

type UpdateCouponRequest struct {
	Code string `json:"code" validate:"required"`
}

type CouponResponse struct {
	ID        uint   `json:"id"`
	Code      string `json:"code"`
	Active    bool   `json:"active"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
