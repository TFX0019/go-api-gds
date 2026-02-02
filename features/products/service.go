package products

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/TFX0019/api-go-gds/features/plans"
	"github.com/google/uuid"
)

type Service interface {
	Create(userID string, req CreateProductRequest) (*ProductResponse, error)
	GetAll(page, limit int) (*PaginatedResponse, error)
	GetByID(id string) (*ProductResponse, error)
	GetByUserID(userID string, page, limit int) (*PaginatedResponse, error)
	Update(id string, req UpdateProductRequest) (*ProductResponse, error)
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

	plan, err := s.plansRepo.FindByProductID(user.Subscription.ProductID)
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
	// For numeric values, need to know if they were provided (e.g. pointer in request or check non-zero?)
	// But struct is values. Assuming full replace or only if non-zero?
	// The prompt implies "update", typically existing values are kept if not provided.
	// However, numbers can be 0.
	// A simple approach is to always update if we trust the request has current/new values,
	// or use pointers in UpdateRequest.
	// Given the prompt didn't specify partial updates intricately, I'll update all fields from request
	// assuming the UI sends the full object state or specific fields?
	// Actually, usually PUT sends full resource, PATCH sends partial.
	// The request is 'UpdateProductRequest' with value types.
	// I'll update them blindly as typical Update logic often does with simple DTOs.
	// IF the user sends 0, it becomes 0.

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

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func mapToResponse(p Product) ProductResponse {
	var cid *string
	if p.ClientID != nil {
		s := p.ClientID.String()
		cid = &s
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
		CreatedAt:            p.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            p.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
