package helps

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(help *Help) error
	FindAll() ([]Help, error)
	FindByTag(tag string) (*Help, error)
	FindByID(id string) (*Help, error)
	Update(help *Help) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(help *Help) error {
	return r.db.Create(help).Error
}

func (r *repository) FindAll() ([]Help, error) {
	var helps []Help
	if err := r.db.Order("created_at desc").Find(&helps).Error; err != nil {
		return nil, err
	}
	return helps, nil
}

func (r *repository) FindByTag(tag string) (*Help, error) {
	var help Help
	if err := r.db.Where("tag = ?", tag).First(&help).Error; err != nil {
		return nil, err
	}
	return &help, nil
}

func (r *repository) FindByID(id string) (*Help, error) {
	var help Help
	if err := r.db.First(&help, id).Error; err != nil {
		return nil, err
	}
	return &help, nil
}

func (r *repository) Update(help *Help) error {
	return r.db.Save(help).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Help{}, id).Error
}
