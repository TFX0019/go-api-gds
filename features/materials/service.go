package materials

import (
	"errors"
	"fmt"
	"strconv"
)

type Service interface {
	Create(userID string, req CreateMaterialRequest, imageURL string) (*MaterialResponse, error)
	GetAll(page, limit int) (*PaginatedResponse, error)
	GetByID(id string) (*MaterialResponse, error)
	GetByUserID(userID string, page, limit int) (*PaginatedResponse, error)
	Update(id string, req UpdateMaterialRequest, imageURL string) (*MaterialResponse, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(userID string, req CreateMaterialRequest, imageURL string) (*MaterialResponse, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	material := &Material{
		UserID:   uint(uid),
		Name:     req.Name,
		Price:    req.Price,
		Quantity: req.Quantity,
		Unit:     req.Unit,
		ImageURL: imageURL,
	}

	if err := s.repo.Create(material); err != nil {
		return nil, err
	}

	res := mapToResponse(*material)
	return &res, nil
}

func (s *service) GetAll(page, limit int) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	materials, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []MaterialResponse
	for _, m := range materials {
		responses = append(responses, mapToResponse(m))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetByID(id string) (*MaterialResponse, error) {
	material, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := mapToResponse(*material)
	return &res, nil
}

func (s *service) GetByUserID(userID string, page, limit int) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	materials, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var responses []MaterialResponse
	for _, m := range materials {
		responses = append(responses, mapToResponse(m))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) Update(id string, req UpdateMaterialRequest, imageURL string) (*MaterialResponse, error) {
	material, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		material.Name = req.Name
	}
	// For numeric/required fields being updated, if they are sent as 0 it might mean "no change" or "set to 0".
	// Since I'm using a simple update struct without pointers for primitives in the DTO (except for the ones I wrote),
	// I'll assume if it's there, it's an update.
	// But in CreateMaterialRequest 'struct tags' are 'form'. in Update 'form'.
	// In the previous step I used value types.
	// Logic: If the user sends the field, Fiber parser fills it. If not, it's 0/empty.
	// A robust update usually separates "set to 0" vs "not provided".
	// For this task, I'll update if non-zero or just update blindly if that's the established pattern.
	// In customers service I mapped everything.
	// I'll map everything here too, as usually form-data updates send the whole object or the client handles it.
	// Note: If Unit is empty string, should I update it? Probably effectively "not provided".

	if req.Name != "" {
		material.Name = req.Name
	}
	if req.Unit != "" {
		material.Unit = req.Unit
	}
	// Price is tricky if 0 is valid. I'll update it for now.
	material.Price = req.Price
	material.Quantity = req.Quantity

	if imageURL != "" {
		material.ImageURL = imageURL
	}

	if err := s.repo.Update(material); err != nil {
		return nil, err
	}

	res := mapToResponse(*material)
	return &res, nil
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func mapToResponse(m Material) MaterialResponse {
	return MaterialResponse{
		ID:        m.ID.String(),
		UserID:    fmt.Sprintf("%d", m.UserID),
		Name:      m.Name,
		Price:     m.Price,
		Quantity:  m.Quantity,
		Unit:      m.Unit,
		ImageURL:  m.ImageURL,
		CreatedAt: m.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: m.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
