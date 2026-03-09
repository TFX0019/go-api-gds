package banners

import (
	"errors"
)

type Service interface {
	Create(image string) (*BannerResponse, error)
	GetAllAdmin(page, limit int) (*PaginatedBannerResponse, error)
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(image string) (*BannerResponse, error) {
	banner := &Banner{
		Image: image,
	}

	if err := s.repo.Create(banner); err != nil {
		return nil, err
	}

	res := mapToResponse(*banner)
	return &res, nil
}

func (s *service) GetAllAdmin(page, limit int) (*PaginatedBannerResponse, error) {
	offset := (page - 1) * limit
	banners, total, err := s.repo.FindAll(limit, offset)
	if err != nil {
		return nil, err
	}

	var data []BannerResponse
	for _, b := range banners {
		data = append(data, mapToResponse(b))
	}

	return &PaginatedBannerResponse{
		Data:  data,
		Total: total,
		Page:  page,
		Limit: limit,
	}, nil
}

func (s *service) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("banner not found")
	}
	return s.repo.Delete(id)
}

func mapToResponse(b Banner) BannerResponse {
	return BannerResponse{
		ID:        b.ID.String(),
		Image:     b.Image,
		CreatedAt: b.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: b.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
