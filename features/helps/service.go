package helps

import (
	"errors"
)

type Service interface {
	Create(req CreateHelpRequest) (*HelpResponse, error)
	GetAll() ([]HelpResponse, error)
	GetByTag(tag string) (*HelpResponse, error)
	Update(id string, req UpdateHelpRequest) (*HelpResponse, error)
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

func (s *service) Create(req CreateHelpRequest) (*HelpResponse, error) {
	// Check if tag already exists
	_, err := s.repo.FindByTag(req.Tag)
	if err == nil {
		return nil, errors.New("tag already exists")
	}

	help := &Help{
		Tag:         req.Tag,
		Description: req.Description,
		Active:      true,
	}

	if err := s.repo.Create(help); err != nil {
		return nil, err
	}

	res := mapToResponse(*help)
	return &res, nil
}

func (s *service) GetAll() ([]HelpResponse, error) {
	helps, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []HelpResponse
	for _, h := range helps {
		res = append(res, mapToResponse(h))
	}
	return res, nil
}

func (s *service) GetByTag(tag string) (*HelpResponse, error) {
	help, err := s.repo.FindByTag(tag)
	if err != nil {
		return nil, errors.New("help not found")
	}

	res := mapToResponse(*help)
	return &res, nil
}

func (s *service) Update(id string, req UpdateHelpRequest) (*HelpResponse, error) {
	help, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("help not found")
	}

	if req.Tag != "" {
		if req.Tag != help.Tag {
			_, errTag := s.repo.FindByTag(req.Tag)
			if errTag == nil {
				return nil, errors.New("tag already exists")
			}
		}
		help.Tag = req.Tag
	}
	if req.Description != "" {
		help.Description = req.Description
	}

	if err := s.repo.Update(help); err != nil {
		return nil, err
	}

	res := mapToResponse(*help)
	return &res, nil
}

func (s *service) Activate(id string) error {
	help, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("help not found")
	}
	help.Active = true
	return s.repo.Update(help)
}

func (s *service) Deactivate(id string) error {
	help, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("help not found")
	}
	help.Active = false
	return s.repo.Update(help)
}

func (s *service) Delete(id string) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("help not found")
	}
	return s.repo.Delete(id)
}

func mapToResponse(h Help) HelpResponse {
	return HelpResponse{
		ID:          h.ID,
		Tag:         h.Tag,
		Description: h.Description,
		Active:      h.Active,
		CreatedAt:   h.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   h.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
