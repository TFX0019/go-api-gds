package links

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(link *Link) error
	FindAll() ([]Link, error)
	FindByID(id string) (*Link, error)
	Update(link *Link) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(link *Link) error {
	return r.db.Create(link).Error
}

func (r *repository) FindAll() ([]Link, error) {
	var links []Link
	if err := r.db.Order("created_at desc").Find(&links).Error; err != nil {
		return nil, err
	}
	return links, nil
}

func (r *repository) FindByID(id string) (*Link, error) {
	var link Link
	if err := r.db.First(&link, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &link, nil
}

func (r *repository) Update(link *Link) error {
	return r.db.Save(link).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Link{}, "id = ?", id).Error
}
