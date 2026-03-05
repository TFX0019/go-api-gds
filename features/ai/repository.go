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

	// AISuggestions
	CreateSuggestion(suggestion *AISuggestion) error
	UpdateSuggestion(suggestion *AISuggestion) error
	DeleteSuggestion(id string) error
	FindSuggestionByID(id string) (*AISuggestion, error)
	FindAllSuggestions() ([]AISuggestion, error)
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

func (r *repository) CreateSuggestion(suggestion *AISuggestion) error {
	return r.db.Create(suggestion).Error
}

func (r *repository) UpdateSuggestion(suggestion *AISuggestion) error {
	return r.db.Save(suggestion).Error
}

func (r *repository) DeleteSuggestion(id string) error {
	return r.db.Delete(&AISuggestion{}, "id = ?", id).Error
}

func (r *repository) FindSuggestionByID(id string) (*AISuggestion, error) {
	var sug AISuggestion
	if err := r.db.First(&sug, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sug, nil
}

func (r *repository) FindAllSuggestions() ([]AISuggestion, error) {
	var sugs []AISuggestion
	if err := r.db.Order("created_at desc").Find(&sugs).Error; err != nil {
		return nil, err
	}
	return sugs, nil
}
