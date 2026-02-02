package plans

import "gorm.io/gorm"

type Repository interface {
	FindByProductID(productID string) (*Plan, error)
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
