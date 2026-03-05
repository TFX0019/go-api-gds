package support

import "gorm.io/gorm"

type Repository interface {
	// Category ops
	CreateCategory(req *SupportCategory) error
	FindAllCategories() ([]SupportCategory, error)
	FindCategoryByID(id string) (*SupportCategory, error)
	UpdateCategory(req *SupportCategory) error

	// Support ops
	CreateSupport(req *Support) error
	FindSupportByID(id string) (*Support, error)
	FindAllParentSupports(limit, offset int) ([]Support, int64, error)
	FindByUserID(userID uint, limit, offset int) ([]Support, int64, error)
	FindRepliesByParentID(parentID string) ([]Support, error)
	UpdateSupport(req *Support) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateCategory(req *SupportCategory) error {
	return r.db.Create(req).Error
}

func (r *repository) FindAllCategories() ([]SupportCategory, error) {
	var categories []SupportCategory
	err := r.db.Order("created_at desc").Find(&categories).Error
	return categories, err
}

func (r *repository) FindCategoryByID(id string) (*SupportCategory, error) {
	var cat SupportCategory
	err := r.db.Where("id = ?", id).First(&cat).Error
	if err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *repository) UpdateCategory(req *SupportCategory) error {
	return r.db.Save(req).Error
}

func (r *repository) CreateSupport(req *Support) error {
	return r.db.Create(req).Error
}

func (r *repository) FindSupportByID(id string) (*Support, error) {
	var support Support
	err := r.db.Preload("SupportCategory").Where("id = ?", id).First(&support).Error
	if err != nil {
		return nil, err
	}
	return &support, nil
}

func (r *repository) FindAllParentSupports(limit, offset int) ([]Support, int64, error) {
	var supports []Support
	var total int64

	query := r.db.Model(&Support{}).Where("parent_id IS NULL")

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("SupportCategory").Limit(limit).Offset(offset).Order("created_at desc").Find(&supports).Error
	if err != nil {
		return nil, 0, err
	}

	return supports, total, nil
}

func (r *repository) FindByUserID(userID uint, limit, offset int) ([]Support, int64, error) {
	var supports []Support
	var total int64

	query := r.db.Model(&Support{}).Where("user_id = ? AND is_deleted = false AND parent_id IS NULL", userID)

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Preload("SupportCategory").Limit(limit).Offset(offset).Order("created_at desc").Find(&supports).Error
	if err != nil {
		return nil, 0, err
	}

	return supports, total, nil
}

func (r *repository) FindRepliesByParentID(parentID string) ([]Support, error) {
	var supports []Support
	err := r.db.Preload("SupportCategory").Where("parent_id = ?", parentID).Order("created_at asc").Find(&supports).Error
	return supports, err
}

func (r *repository) UpdateSupport(req *Support) error {
	return r.db.Save(req).Error
}
