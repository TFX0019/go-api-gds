package user

import (
	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/TFX0019/api-go-gds/pkg/utils"
	"gorm.io/gorm"
)

type Repository interface {
	FindAllWithPlan(pagination *utils.Pagination) (*utils.Pagination, error)
	FindByID(id uint) (*auth.User, error)
	UpdateStatus(id uint, isActive bool) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAllWithPlan(pagination *utils.Pagination) (*utils.Pagination, error) {
	var users []UserListDTO
	var totalRows int64

	// Count total rows
	r.db.Table("users").Count(&totalRows)
	pagination.TotalRows = totalRows

	totalPages := int(totalRows/int64(pagination.GetLimit())) + 1
	if totalRows%int64(pagination.GetLimit()) == 0 {
		totalPages = int(totalRows / int64(pagination.GetLimit()))
	}
	pagination.TotalPages = totalPages

	// Query with pagination
	err := r.db.Table("users").
		Select("users.id, users.name, users.email, users.is_active, users.is_verified, plans.title as plan_name, subscriptions.status as subscription_status").
		Joins("LEFT JOIN subscriptions ON subscriptions.user_id = users.id").
		Joins("LEFT JOIN plans ON plans.product_id = subscriptions.product_id").
		Limit(pagination.GetLimit()).
		Offset(pagination.GetOffset()).
		Scan(&users).Error

	if err != nil {
		return nil, err
	}
	pagination.Rows = users
	return pagination, nil
}

func (r *repository) FindByID(id uint) (*auth.User, error) {
	var user auth.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateStatus(id uint, isActive bool) error {
	return r.db.Model(&auth.User{}).Where("id = ?", id).Update("is_active", isActive).Error
}
