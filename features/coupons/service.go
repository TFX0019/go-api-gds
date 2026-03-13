package coupons

import (
	"errors"
)

type Service interface {
	Get() (*CouponResponse, error)
	Update(req UpdateCouponRequest) (*CouponResponse, error)
	Activate() error
	Deactivate() error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get() (*CouponResponse, error) {
	coupon, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	res := mapToResponse(*coupon)
	return &res, nil
}

func (s *service) Update(req UpdateCouponRequest) (*CouponResponse, error) {
	coupon, err := s.repo.Get()
	if err != nil {
		return nil, errors.New("failed to retrieve current coupon")
	}

	coupon.Code = req.Code

	if err := s.repo.Update(coupon); err != nil {
		return nil, err
	}

	res := mapToResponse(*coupon)
	return &res, nil
}

func (s *service) Activate() error {
	coupon, err := s.repo.Get()
	if err != nil {
		return errors.New("failed to retrieve current coupon")
	}
	coupon.Active = true
	return s.repo.Update(coupon)
}

func (s *service) Deactivate() error {
	coupon, err := s.repo.Get()
	if err != nil {
		return errors.New("failed to retrieve current coupon")
	}
	coupon.Active = false
	return s.repo.Update(coupon)
}

func mapToResponse(c Coupon) CouponResponse {
	return CouponResponse{
		ID:        c.ID,
		Code:      c.Code,
		Active:    c.Active,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
