package materials

import "gorm.io/gorm"

type Repository interface {
	Create(material *Material) error
	FindAll(limit, offset int) ([]Material, int64, error)
	FindByID(id string) (*Material, error)
	FindByUserID(userID string, limit, offset int) ([]Material, int64, error)
	Update(material *Material) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(material *Material) error {
	return r.db.Create(material).Error
}

func (r *repository) FindAll(limit, offset int) ([]Material, int64, error) {
	var materials []Material
	var total int64

	err := r.db.Model(&Material{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&materials).Error
	if err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

func (r *repository) FindByID(id string) (*Material, error) {
	var material Material
	err := r.db.Where("id = ?", id).First(&material).Error
	if err != nil {
		return nil, err
	}
	return &material, nil
}

func (r *repository) FindByUserID(userID string, limit, offset int) ([]Material, int64, error) {
	var materials []Material
	var total int64

	err := r.db.Model(&Material{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at desc").Find(&materials).Error
	if err != nil {
		return nil, 0, err
	}

	return materials, total, nil
}

func (r *repository) Update(material *Material) error {
	return r.db.Save(material).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Material{}, "id = ?", id).Error
}
