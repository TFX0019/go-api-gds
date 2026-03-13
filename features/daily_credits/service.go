package daily_credits

import (
	"errors"
)

type Service interface {
	Get() (*DailyCreditResponse, error)
	Update(req UpdateDailyCreditRequest) (*DailyCreditResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Get() (*DailyCreditResponse, error) {
	credit, err := s.repo.Get()
	if err != nil {
		return nil, err
	}

	res := mapToResponse(*credit)
	return &res, nil
}

func (s *service) Update(req UpdateDailyCreditRequest) (*DailyCreditResponse, error) {
	credit, err := s.repo.Get()
	if err != nil {
		return nil, errors.New("failed to retrieve current configuration")
	}

	credit.Free = req.Free
	credit.Premium = req.Premium

	if err := s.repo.Update(credit); err != nil {
		return nil, err
	}

	res := mapToResponse(*credit)
	return &res, nil
}

func mapToResponse(c DailyCredit) DailyCreditResponse {
	return DailyCreditResponse{
		ID:        c.ID,
		Free:      c.Free,
		Premium:   c.Premium,
		CreatedAt: c.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: c.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
