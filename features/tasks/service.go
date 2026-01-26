package tasks

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(userID string, req CreateTaskRequest) (*TaskResponse, error)
	GetAll(page, limit int, status, date string) (*PaginatedResponse, error)
	GetByID(id string) (*TaskResponse, error)
	GetByUserID(userID string, page, limit int, status, date string) (*PaginatedResponse, error)
	Update(id string, req UpdateTaskRequest) (*TaskResponse, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(userID string, req CreateTaskRequest) (*TaskResponse, error) {
	uid, err := strconv.ParseUint(userID, 10, 32)
	if err != nil {
		return nil, errors.New("invalid user id")
	}

	dateTime, err := time.Parse(time.RFC3339, req.DateTime)
	if err != nil {
		return nil, errors.New("invalid date_time format")
	}

	var prodUUID *uuid.UUID
	if req.ProductID != nil && *req.ProductID != "" {
		id, err := uuid.Parse(*req.ProductID)
		if err != nil {
			return nil, errors.New("invalid product id")
		}
		prodUUID = &id
	}

	task := &Task{
		UserID:      uint(uid),
		Name:        req.Name,
		Description: req.Description,
		Status:      req.Status,
		DateTime:    dateTime,
		ProductID:   prodUUID,
	}

	if err := s.repo.Create(task); err != nil {
		return nil, err
	}

	// Re-fetch to populate associations (Product, Client) for the response
	createdTask, err := s.repo.FindByID(task.ID.String())
	if err == nil {
		res := mapToResponse(*createdTask)
		return &res, nil
	}

	res := mapToResponse(*task)
	return &res, nil
}

func (s *service) GetAll(page, limit int, status, date string) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	tasks, total, err := s.repo.FindAll(limit, offset, status, date)
	if err != nil {
		return nil, err
	}

	var responses []TaskResponse
	for _, t := range tasks {
		responses = append(responses, mapToResponse(t))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetByID(id string) (*TaskResponse, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	res := mapToResponse(*task)
	return &res, nil
}

func (s *service) GetByUserID(userID string, page, limit int, status, date string) (*PaginatedResponse, error) {
	offset := (page - 1) * limit
	tasks, total, err := s.repo.FindByUserID(userID, limit, offset, status, date)
	if err != nil {
		return nil, err
	}

	var responses []TaskResponse
	for _, t := range tasks {
		responses = append(responses, mapToResponse(t))
	}

	return &PaginatedResponse{
		Data:  responses,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) Update(id string, req UpdateTaskRequest) (*TaskResponse, error) {
	task, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	if req.Name != "" {
		task.Name = req.Name
	}
	if req.Description != "" {
		task.Description = req.Description
	}
	if req.Status != "" {
		task.Status = req.Status
	}
	if req.DateTime != "" {
		dt, err := time.Parse(time.RFC3339, req.DateTime)
		if err == nil {
			task.DateTime = dt
		}
	}

	if req.ProductID != nil {
		if *req.ProductID == "" {
			task.ProductID = nil
		} else {
			pid, err := uuid.Parse(*req.ProductID)
			if err != nil {
				return nil, errors.New("invalid product id")
			}
			task.ProductID = &pid
		}
	}

	if err := s.repo.Update(task); err != nil {
		return nil, err
	}

	// Re-fetch to get updated associations or just return what we have (associations present from FindByID)
	// If ProductID changed, FindByID data is stale for Product.
	// Ideally we re-fetch.
	updatedTask, err := s.repo.FindByID(id)
	if err == nil {
		res := mapToResponse(*updatedTask)
		return &res, nil
	}

	res := mapToResponse(*task)
	return &res, nil
}

func (s *service) Delete(id string) error {
	return s.repo.Delete(id)
}

func mapToResponse(t Task) TaskResponse {
	var prodInfo *ProductInfo

	if t.Product != nil {
		info := &ProductInfo{
			ID:   t.Product.ID.String(),
			Name: t.Product.Name,
		}

		if t.Product.Client != nil {
			info.Client = &ClientInfo{
				Name:      t.Product.Client.Name,
				Phone:     t.Product.Client.Phone,
				Email:     t.Product.Client.Email,
				AvatarURL: t.Product.Client.AvatarURL,
			}
		}
		prodInfo = info
	}

	return TaskResponse{
		ID:          t.ID.String(),
		UserID:      fmt.Sprintf("%d", t.UserID),
		Name:        t.Name,
		Description: t.Description,
		Status:      t.Status,
		DateTime:    t.DateTime.Format(time.RFC3339),
		Product:     prodInfo,
		CreatedAt:   t.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   t.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
