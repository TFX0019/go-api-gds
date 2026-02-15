package user

import (
	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/TFX0019/api-go-gds/pkg/utils"
)

type Service interface {
	ListUsers(pagination *utils.Pagination) (*utils.Pagination, error)
	ActivateUser(id uint) (*auth.User, error)
	BanUser(id uint) (*auth.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) ListUsers(pagination *utils.Pagination) (*utils.Pagination, error) {
	return s.repo.FindAllWithPlan(pagination)
}

func (s *service) ActivateUser(id uint) (*auth.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if user.IsActive {
		return user, nil // Already active
	}
	if err := s.repo.UpdateStatus(id, true); err != nil {
		return nil, err
	}
	user.IsActive = true
	return user, nil
}

func (s *service) BanUser(id uint) (*auth.User, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if !user.IsActive {
		return user, nil // Already banned/inactive
	}
	if err := s.repo.UpdateStatus(id, false); err != nil {
		return nil, err
	}
	user.IsActive = false
	return user, nil
}
