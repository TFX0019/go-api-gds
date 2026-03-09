package links

import (
	"errors"
)

type Service interface {
	Create(req CreateLinkRequest) (*LinkResponse, error)
	GetAll() ([]LinkResponse, error)
	Update(id string, req UpdateLinkRequest) (*LinkResponse, error)
	Activate(id string) error
	Deactivate(id string) error
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(req CreateLinkRequest) (*LinkResponse, error) {
	link := &Link{
		Title:       req.Title,
		Description: req.Description,
		URL:         req.URL,
		Active:      true,
	}

	if err := s.repo.Create(link); err != nil {
		return nil, err
	}

	res := mapToResponse(*link)
	return &res, nil
}

func (s *service) GetAll() ([]LinkResponse, error) {
	links, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []LinkResponse
	for _, l := range links {
		res = append(res, mapToResponse(l))
	}
	return res, nil
}

func (s *service) Update(id string, req UpdateLinkRequest) (*LinkResponse, error) {
	link, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("link not found")
	}

	if req.Title != "" {
		link.Title = req.Title
	}
	if req.Description != "" {
		link.Description = req.Description
	}
	if req.URL != "" {
		link.URL = req.URL
	}

	if err := s.repo.Update(link); err != nil {
		return nil, err
	}

	res := mapToResponse(*link)
	return &res, nil
}

func (s *service) Activate(id string) error {
	link, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("link not found")
	}
	link.Active = true
	return s.repo.Update(link)
}

func (s *service) Deactivate(id string) error {
	link, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("link not found")
	}
	link.Active = false
	return s.repo.Update(link)
}

func (s *service) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("link not found")
	}
	return s.repo.Delete(id)
}

func mapToResponse(l Link) LinkResponse {
	return LinkResponse{
		ID:          l.ID.String(),
		Title:       l.Title,
		Description: l.Description,
		URL:         l.URL,
		Active:      l.Active,
		CreatedAt:   l.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   l.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
