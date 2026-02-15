package plans

import (
	"errors"

	"gorm.io/gorm"
)

type Repository interface {
	FindByProductID(productID string) (*Plan, error)
	FindAll() ([]Plan, error)
	FindAllActive() ([]Plan, error)
	FindByID(id uint) (*Plan, error)
	Update(plan *Plan) error
	Create(plan *Plan) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByProductID(productID string) (*Plan, error) {
	var plan Plan
	err := r.db.Where("product_id = ?", productID).First(&plan).Error
	if err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *repository) FindAll() ([]Plan, error) {
	var plans []Plan
	err := r.db.Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *repository) FindAllActive() ([]Plan, error) {
	var plans []Plan
	err := r.db.Where("is_active = ?", true).Find(&plans).Error
	if err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *repository) FindByID(id uint) (*Plan, error) {
	var plan Plan
	err := r.db.First(&plan, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("plan not found")
		}
		return nil, err
	}
	return &plan, nil
}

func (r *repository) Update(plan *Plan) error {
	return r.db.Save(plan).Error
}

func (r *repository) Create(plan *Plan) error {
	return r.db.Create(plan).Error
}
