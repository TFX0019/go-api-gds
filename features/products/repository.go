package products

import "gorm.io/gorm"

type Repository interface {
	Create(product *Product) error
	FindAll(limit, offset int) ([]Product, int64, error)
	FindByID(id string) (*Product, error)
	FindByUserID(userID string, limit, offset int) ([]Product, int64, error)
	CountByUserID(userID uint) (int64, error)
	GetProfitLoss(userID string, month int) (*ProfitLossResponse, error)
	Update(product *Product) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(product *Product) error {
	return r.db.Create(product).Error
}

func (r *repository) FindAll(limit, offset int) ([]Product, int64, error) {
	var products []Product
	var total int64

	err := r.db.Model(&Product{}).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *repository) FindByID(id string) (*Product, error) {
	var product Product
	err := r.db.Where("id = ?", id).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *repository) FindByUserID(userID string, limit, offset int) ([]Product, int64, error) {
	var products []Product
	var total int64

	err := r.db.Model(&Product{}).Where("user_id = ?", userID).Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = r.db.Where("user_id = ?", userID).Limit(limit).Offset(offset).Order("created_at desc").Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *repository) CountByUserID(userID uint) (int64, error) {
	var total int64
	err := r.db.Model(&Product{}).Where("user_id = ?", userID).Count(&total).Error
	return total, err
}

func (r *repository) GetProfitLoss(userID string, month int) (*ProfitLossResponse, error) {
	var result ProfitLossResponse
	query := r.db.Model(&Product{}).Where("user_id = ? AND status = ?", userID, "paid")

	if month > 0 {
		query = query.Where("EXTRACT(MONTH FROM date_paid) = ? AND EXTRACT(YEAR FROM date_paid) = EXTRACT(YEAR FROM CURRENT_DATE)", month)
	}

	err := query.Select("COALESCE(SUM(materials_cost), 0) as total_materials_cost, COALESCE(SUM(hours_cost), 0) as total_hours_cost, COALESCE(SUM(fixed_expenses_amount), 0) as total_fixed_expenses_amount, COALESCE(SUM(profit_amount), 0) as total_profit_amount").
		Scan(&result).Error
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (r *repository) Update(product *Product) error {
	return r.db.Save(product).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Product{}, "id = ?", id).Error
}
