package customers

import (
	"gorm.io/gorm"
)

type Repository interface {
	Create(customer *Customer) error
	FindAll(limit, offset int) ([]Customer, int64, error)
	FindByID(id string) (*Customer, error)
	FindByUserID(userID string, limit, offset int) ([]Customer, int64, error)
	CountByUserID(userID uint) (int64, error)
	Update(customer *Customer) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(customer *Customer) error {
	return r.db.Create(customer).Error
}

func (r *repository) FindAll(limit, offset int) ([]Customer, int64, error) {
	var customers []Customer
	var total int64

	err := r.db.Model(&Customer{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *repository) FindByID(id string) (*Customer, error) {
	var customer Customer
	err := r.db.Where("id = ?", id).First(&customer).Error
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (r *repository) FindByUserID(userID string, limit, offset int) ([]Customer, int64, error) {
	var customers []Customer
	var total int64

	err := r.db.Model(&Customer{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at desc").Find(&customers).Error
	if err != nil {
		return nil, 0, err
	}

	return customers, total, nil
}

func (r *repository) CountByUserID(userID uint) (int64, error) {
	var total int64
	err := r.db.Model(&Customer{}).Where("user_id = ?", userID).Count(&total).Error
	return total, err
}

func (r *repository) Update(customer *Customer) error {
	return r.db.Save(customer).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Customer{}, "id = ?", id).Error
}
