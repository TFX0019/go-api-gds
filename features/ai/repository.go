package ai

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(generation *AIGeneration) error
	Update(generation *AIGeneration) error
	FindByID(id string) (*AIGeneration, error)
	FindByUserID(userID uint, limit, offset int) ([]AIGeneration, int64, error)
	FindAll(limit, offset int) ([]AIGeneration, int64, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(generation *AIGeneration) error {
	return r.db.Create(generation).Error
}

func (r *repository) Update(generation *AIGeneration) error {
	return r.db.Save(generation).Error
}

func (r *repository) FindByID(id string) (*AIGeneration, error) {
	var gen AIGeneration
	if err := r.db.Preload("User").First(&gen, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &gen, nil
}

func (r *repository) FindByUserID(userID uint, limit, offset int) ([]AIGeneration, int64, error) {
	var gens []AIGeneration
	var total int64

	query := r.db.Model(&AIGeneration{}).Where("user_id = ?", userID)
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Limit(limit).Offset(offset).Order("created_at desc").Find(&gens).Error; err != nil {
		return nil, 0, err
	}

	return gens, total, nil
}

func (r *repository) FindAll(limit, offset int) ([]AIGeneration, int64, error) {
	var gens []AIGeneration
	var total int64

	query := r.db.Model(&AIGeneration{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Preload("User").Limit(limit).Offset(offset).Order("created_at desc").Find(&gens).Error; err != nil {
		return nil, 0, err
	}

	return gens, total, nil
}
