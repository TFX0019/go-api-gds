package banners

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(banner *Banner) error
	FindAll(limit, offset int) ([]Banner, int64, error)
	FindByID(id string) (*Banner, error)
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(banner *Banner) error {
	return r.db.Create(banner).Error
}

func (r *repository) FindAll(limit, offset int) ([]Banner, int64, error) {
	var banners []Banner
	var total int64

	query := r.db.Model(&Banner{})
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if err := query.Limit(limit).Offset(offset).Order("created_at desc").Find(&banners).Error; err != nil {
		return nil, 0, err
	}

	return banners, total, nil
}

func (r *repository) FindByID(id string) (*Banner, error) {
	var banner Banner
	if err := r.db.First(&banner, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &banner, nil
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Banner{}, "id = ?", id).Error
}
