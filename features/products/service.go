package products

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/TFX0019/api-go-gds/features/plans"
	"github.com/google/uuid"
)

type Service interface {
	Create(userID string, req CreateProductRequest) (*ProductResponse, error)
	GetAll(page, limit int) (*PaginatedResponse, error)
	GetByID(id string) (*ProductResponse, error)
	GetByUserID(userID string, page, limit int) (*PaginatedResponse, error)
	GetProfitLoss(userID string, month int) (*ProfitLossResponse, error)
	Update(id string, req UpdateProductRequest) (*ProductResponse, error)
	UpdateStatus(id string, status string) (*ProductResponse, error)
	Delete(id string) error
}

type service struct {
	repo      Repository
	authRepo  auth.Repository
	plansRepo plans.Repository
}

func NewService(repo Repository, authRepo auth.Repository, plansRepo plans.Repository) Service {
	return &service{repo: repo, authRepo: authRepo, plansRepo: plansRepo}
}

func (s *service) Create(userID string, req CreateProductRequest) (*ProductResponse, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	// Check Limits
	user, err := s.authRepo.FindByID(uint(uid))
	if err != nil {
		return nil, errors.New("user not found")
	}

	productID := user.Subscription.ProductID
	if productID == "" {
		productID = "free_tier"
	}

	plan, err := s.plansRepo.FindByProductID(productID)
	if err != nil {
		return nil, errors.New("plan not found")
	}

	if plan.MaxProducts != -1 {
		count, err := s.repo.CountByUserID(uint(uid))
		if err != nil {
			return nil, err
		}
		if int(count) >= plan.MaxProducts {
			return nil, fmt.Errorf("product limit reached for your plan (%d)", plan.MaxProducts)
		}
	}

	var clientUUID *uuid.UUID
	if req.ClientID != nil && *req.ClientID != "" {
		id, err := uuid.Parse(*req.ClientID)
		if err != nil {
			return nil, errors.New("invalid client id")
		}
		clientUUID = &id
	}

	product := &Product{
		UserID:               uint(uid),
		Name:                 req.Name,
		ClientID:             clientUUID,
		MaterialsCost:        req.MaterialsCost,
		HoursCost:            req.HoursCost,
		ProfitPercentage:     req.ProfitPercentage,
		IncludeFixedExpenses: req.IncludeFixedExpenses,
		FixedExpenseRate:     req.FixedExpenseRate,
		Subtotal:             req.Subtotal,
		FixedExpensesAmount:  req.FixedExpensesAmount,
		BaseTotal:            req.BaseTotal,
		ProfitAmount:         req.ProfitAmount,
		Total:                req.Total,
		Status:               "pending",
	}

	if err := s.repo.Create(product); err != nil {
		return nil, err
	}

	res := mapToResponse(*product)
	return &res, nil
}

func (s *service) GetAll(page, limit int) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	products, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []ProductResponse
	for _, p := range products {
		responses = append(responses, mapToResponse(p))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetByID(id string) (*ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := mapToResponse(*product)
	return &res, nil
}

func (s *service) GetByUserID(userID string, page, limit int) (*PaginatedResponse, error) {
	// Validate userID is numeric since DB expects it? Or handled by repo?
	// Repo uses string type for where clause "user_id = ?", but DB column is int.
	// Postgres driver usually handles string representation of number fine or we can parse.
	// We'll trust GORM/driver unless it fails.

	offset := (page - 1) * limit
	products, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []ProductResponse
	for _, p := range products {
		responses = append(responses, mapToResponse(p))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetProfitLoss(userID string, month int) (*ProfitLossResponse, error) {
	return s.repo.GetProfitLoss(userID, month)
}

func (s *service) Update(id string, req UpdateProductRequest) (*ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.ClientID != nil {
		if *req.ClientID == "" {
			product.ClientID = nil
		} else {
			cid, err := uuid.Parse(*req.ClientID)
			if err != nil {
				return nil, errors.New("invalid client id")
			}
			product.ClientID = &cid
		}
	}

	if req.Name != "" {
		product.Name = req.Name
	}

	product.MaterialsCost = req.MaterialsCost
	product.HoursCost = req.HoursCost
	product.ProfitPercentage = req.ProfitPercentage
	product.IncludeFixedExpenses = req.IncludeFixedExpenses
	product.FixedExpenseRate = req.FixedExpenseRate
	product.Subtotal = req.Subtotal
	product.FixedExpensesAmount = req.FixedExpensesAmount
	product.BaseTotal = req.BaseTotal
	product.ProfitAmount = req.ProfitAmount
	product.Total = req.Total

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	res := mapToResponse(*product)
	return &res, nil
}

func (s *service) UpdateStatus(id string, status string) (*ProductResponse, error) {
	product, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	product.Status = status
	if status == "paid" {
		now := time.Now()
		product.DatePaid = &now
	} else {
		product.DatePaid = nil
	}

	if err := s.repo.Update(product); err != nil {
		return nil, err
	}

	res := mapToResponse(*product)
	return &res, nil
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func mapToResponse(p Product) ProductResponse {
	var cid *string
	if p.ClientID != nil {
		s := p.ClientID.String()
		cid = &s
	}

	var datePaid *string
	if p.DatePaid != nil {
		s := p.DatePaid.Format("2006-01-02 15:04:05")
		datePaid = &s
	}

	return ProductResponse{
		ID:                   p.ID.String(),
		UserID:               fmt.Sprintf("%d", p.UserID),
		Name:                 p.Name,
		ClientID:             cid,
		MaterialsCost:        p.MaterialsCost,
		HoursCost:            p.HoursCost,
		ProfitPercentage:     p.ProfitPercentage,
		IncludeFixedExpenses: p.IncludeFixedExpenses,
		FixedExpenseRate:     p.FixedExpenseRate,
		Subtotal:             p.Subtotal,
		FixedExpensesAmount:  p.FixedExpensesAmount,
		BaseTotal:            p.BaseTotal,
		ProfitAmount:         p.ProfitAmount,
		Total:                p.Total,
		Status:               p.Status,
		DatePaid:             datePaid,
		CreatedAt:            p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
