package customers

import (
	"errors"
	"fmt"
	"strconv"
)

type Service interface {
	Create(userID string, req CreateCustomerRequest, avatarURL string) error
	GetAll(page, limit int) (*PaginatedResponse, error)
	GetByID(id string) (*CustomerResponse, error)
	GetByUserID(userID string, page, limit int) (*PaginatedResponse, error)
	Update(id string, req UpdateCustomerRequest, avatarURL string) error
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(userID string, req CreateCustomerRequest, avatarURL string) error {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return errors.New("invalid user id")
	}

	customer := &Customer{
		UserID:           uint(uid),
		AvatarURL:        avatarURL,
		Name:             req.Name,
		Phone:            req.Phone,
		Email:            req.Email,
		UsesStandardSize: req.UsesStandardSize,
		Back:             req.Back,
		Neck:             req.Neck,
		FrontSize:        req.FrontSize,
		Armhole:          req.Armhole,
		BackSize:         req.BackSize,
		BustChest:        req.BustChest,
		Waist:            req.Waist,
		Hip:              req.Hip,
		RiseHeight:       req.RiseHeight,
		SkirtLength:      req.SkirtLength,
		PantsLength:      req.PantsLength,
		KneeWidth:        req.KneeWidth,
		HemWidth:         req.HemWidth,
		SleeveLength:     req.SleeveLength,
		CuffSize:         req.CuffSize,
	}

	return s.repo.Create(customer)
}

func (s *service) GetAll(page, limit int) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	customers, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []CustomerResponse
	for _, c := range customers {
		responses = append(responses, mapToResponse(c))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetByID(id string) (*CustomerResponse, error) {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := mapToResponse(*customer)
	return &res, nil
}

func (s *service) GetByUserID(userID string, page, limit int) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	customers, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []CustomerResponse
	for _, c := range customers {
		responses = append(responses, mapToResponse(c))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) Update(id string, req UpdateCustomerRequest, avatarURL string) error {
	customer, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}

	// Update fields
	customer.Name = req.Name
	customer.Phone = req.Phone
	customer.Email = req.Email
	customer.UsesStandardSize = req.UsesStandardSize
	customer.Back = req.Back
	customer.Neck = req.Neck
	customer.FrontSize = req.FrontSize
	customer.Armhole = req.Armhole
	customer.BackSize = req.BackSize
	customer.BustChest = req.BustChest
	customer.Waist = req.Waist
	customer.Hip = req.Hip
	customer.RiseHeight = req.RiseHeight
	customer.SkirtLength = req.SkirtLength
	customer.PantsLength = req.PantsLength
	customer.KneeWidth = req.KneeWidth
	customer.HemWidth = req.HemWidth
	customer.SleeveLength = req.SleeveLength
	customer.CuffSize = req.CuffSize

	if avatarURL != "" {
		customer.AvatarURL = avatarURL
	}

	return s.repo.Update(customer)
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func mapToResponse(c Customer) CustomerResponse {
	return CustomerResponse{
		ID:               c.ID.String(),
		UserID:           fmt.Sprintf("%d", c.UserID),
		AvatarURL:        c.AvatarURL,
		Name:             c.Name,
		Phone:            c.Phone,
		Email:            c.Email,
		UsesStandardSize: c.UsesStandardSize,
		Back:             c.Back,
		Neck:             c.Neck,
		FrontSize:        c.FrontSize,
		Armhole:          c.Armhole,
		BackSize:         c.BackSize,
		BustChest:        c.BustChest,
		Waist:            c.Waist,
		Hip:              c.Hip,
		RiseHeight:       c.RiseHeight,
		SkirtLength:      c.SkirtLength,
		PantsLength:      c.PantsLength,
		KneeWidth:        c.KneeWidth,
		HemWidth:         c.HemWidth,
		SleeveLength:     c.SleeveLength,
		CuffSize:         c.CuffSize,
		CreatedAt:        c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
