package support

import (
	"errors"

	"github.com/google/uuid"
)

type Service interface {
	// Category
	CreateCategory(req CreateCategoryRequest) (*SupportCategoryResponse, error)
	GetAllCategories() ([]SupportCategoryResponse, error)
	UpdateCategory(id string, req UpdateCategoryRequest) (*SupportCategoryResponse, error)
	ActivateCategory(id string) error
	DeactivateCategory(id string) error

	// Support
	CreateSupport(userID uint, req CreateSupportRequest, imageURL string) (*SupportResponse, error)
	GetParentSupportsByUserID(userID uint, page, limit int) (*PaginatedSupportResponse, error)
	GetAllParentSupportsAdmin(page, limit int) (*PaginatedSupportResponse, error)
	GetRepliesByParentID(parentID string) ([]SupportResponse, error)
	UpdateSupport(id string, req UpdateSupportRequest) (*SupportResponse, error)
	DeleteSupport(id string, reqUserID uint, isAdmin bool) error // Admin hard deletes or skips, User soft deletes
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) CreateCategory(req CreateCategoryRequest) (*SupportCategoryResponse, error) {
	cat := &SupportCategory{
		Title:       req.Title,
		Description: req.Description,
		Active:      true,
	}

	if err := s.repo.CreateCategory(cat); err != nil {
		return nil, err
	}
	res := mapCategoryToResponse(*cat)
	return &res, nil
}

func (s *service) GetAllCategories() ([]SupportCategoryResponse, error) {
	cats, err := s.repo.FindAllCategories()
	if err != nil {
		return nil, err
	}
	var res []SupportCategoryResponse
	for _, c := range cats {
		res = append(res, mapCategoryToResponse(c))
	}
	return res, nil
}

func (s *service) UpdateCategory(id string, req UpdateCategoryRequest) (*SupportCategoryResponse, error) {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return nil, err
	}

	if req.Title != "" {
		cat.Title = req.Title
	}
	if req.Description != "" {
		cat.Description = req.Description
	}
	if req.Active != nil {
		cat.Active = *req.Active
	}

	if err := s.repo.UpdateCategory(cat); err != nil {
		return nil, err
	}

	res := mapCategoryToResponse(*cat)
	return &res, nil
}

func (s *service) ActivateCategory(id string) error {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return err
	}
	cat.Active = true
	return s.repo.UpdateCategory(cat)
}

func (s *service) DeactivateCategory(id string) error {
	cat, err := s.repo.FindCategoryByID(id)
	if err != nil {
		return err
	}
	cat.Active = false
	return s.repo.UpdateCategory(cat)
}

func (s *service) CreateSupport(userID uint, req CreateSupportRequest, imageURL string) (*SupportResponse, error) {
	catID, err := uuid.Parse(req.SupportCategoryID)
	if err != nil {
		return nil, errors.New("invalid category id")
	}

	var parentID *uuid.UUID
	if req.ParentID != nil && *req.ParentID != "" {
		pid, err := uuid.Parse(*req.ParentID)
		if err != nil {
			return nil, errors.New("invalid parent id")
		}
		parentID = &pid
	}

	support := &Support{
		Subject:           req.Subject,
		Description:       req.Description,
		UserID:            userID,
		SupportCategoryID: catID,
		Status:            "open",
		Image:             imageURL,
		ParentID:          parentID,
		IsDeleted:         false,
	}

	if err := s.repo.CreateSupport(support); err != nil {
		return nil, err
	}

	// Fetch category to include in response
	if cat, err := s.repo.FindCategoryByID(catID.String()); err == nil && cat != nil {
		support.SupportCategory = *cat
	}

	res := mapSupportToResponse(*support)
	return &res, nil
}

func (s *service) GetParentSupportsByUserID(userID uint, page, limit int) (*PaginatedSupportResponse, error) {
	offset := (page - 1) * limit
	supports, total, err := s.repo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, err
	}

	var res []SupportResponse
	for _, spt := range supports {
		res = append(res, mapSupportToResponse(spt))
	}

	return &PaginatedSupportResponse{
		Data:  res,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetAllParentSupportsAdmin(page, limit int) (*PaginatedSupportResponse, error) {
	offset := (page - 1) * limit
	supports, total, err := s.repo.FindAllParentSupports(limit, offset)
	if err != nil {
		return nil, err
	}

	var res []SupportResponse
	for _, spt := range supports {
		res = append(res, mapSupportToResponse(spt))
	}

	return &PaginatedSupportResponse{
		Data:  res,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) GetRepliesByParentID(parentID string) ([]SupportResponse, error) {
	supports, err := s.repo.FindRepliesByParentID(parentID)
	if err != nil {
		return nil, err
	}
	var res []SupportResponse
	for _, spt := range supports {
		res = append(res, mapSupportToResponse(spt))
	}
	return res, nil
}

func (s *service) UpdateSupport(id string, req UpdateSupportRequest) (*SupportResponse, error) {
	support, err := s.repo.FindSupportByID(id)
	if err != nil {
		return nil, err
	}

	if req.Subject != "" {
		support.Subject = req.Subject
	}
	if req.Description != "" {
		support.Description = req.Description
	}
	if req.SupportCategoryID != nil && *req.SupportCategoryID != "" {
		catID, err := uuid.Parse(*req.SupportCategoryID)
		if err == nil {
			support.SupportCategoryID = catID
		}
	}
	if req.Status != nil && *req.Status != "" {
		support.Status = *req.Status
	}

	if err := s.repo.UpdateSupport(support); err != nil {
		return nil, err
	}

	res := mapSupportToResponse(*support)
	return &res, nil
}

func (s *service) DeleteSupport(id string, reqUserID uint, isAdmin bool) error {
	support, err := s.repo.FindSupportByID(id)
	if err != nil {
		return err
	}

	if isAdmin {
		support.IsDeleted = true
	} else {
		if support.UserID != reqUserID {
			return errors.New("unauthorized to delete this resource")
		}
		support.IsDeleted = true
	}

	return s.repo.UpdateSupport(support) // Soft delete
}

func mapCategoryToResponse(c SupportCategory) SupportCategoryResponse {
	return SupportCategoryResponse{
		ID:          c.ID.String(),
		Title:       c.Title,
		Description: c.Description,
		Active:      c.Active,
		CreatedAt:   c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}

func mapSupportToResponse(s Support) SupportResponse {
	var catRes *SupportCategoryResponse
	if s.SupportCategory.ID != uuid.Nil {
		res := mapCategoryToResponse(s.SupportCategory)
		catRes = &res
	}

	var pid *string
	if s.ParentID != nil {
		str := s.ParentID.String()
		pid = &str
	}

	return SupportResponse{
		ID:                s.ID.String(),
		Subject:           s.Subject,
		Description:       s.Description,
		UserID:            s.UserID,
		SupportCategoryID: s.SupportCategoryID.String(),
		SupportCategory:   catRes,
		Status:            s.Status,
		Image:             s.Image,
		ParentID:          pid,
		IsDeleted:         s.IsDeleted,
		CreatedAt:         s.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:         s.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
