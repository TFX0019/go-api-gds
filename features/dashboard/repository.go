package dashboard

import (
	"github.com/TFX0019/api-go-gds/features/customers"
	"github.com/TFX0019/api-go-gds/features/materials"
	"github.com/TFX0019/api-go-gds/features/products"
	"github.com/TFX0019/api-go-gds/features/tasks"
	"gorm.io/gorm"
)

type Repository interface {
	GetSummaryCounts(userID string) ([]SummaryItem, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetSummaryCounts(userID string) ([]SummaryItem, error) {
	var customerCount int64
	var taskCount int64
	var materialCount int64
	var productCount int64

	// Count customers with user_id
	if err := r.db.Model(&customers.Customer{}).Where("user_id = ?", userID).Count(&customerCount).Error; err != nil {
		return nil, err
	}

	// Count tasks with user_id
	if err := r.db.Model(&tasks.Task{}).Where("user_id = ?", userID).Count(&taskCount).Error; err != nil {
		return nil, err
	}

	// Count materials with user_id
	if err := r.db.Model(&materials.Material{}).Where("user_id = ?", userID).Count(&materialCount).Error; err != nil {
		return nil, err
	}

	// Count products with user_id
	if err := r.db.Model(&products.Product{}).Where("user_id = ?", userID).Count(&productCount).Error; err != nil {
		return nil, err
	}

	return []SummaryItem{
		{Title: "Clientes", Count: customerCount},
		{Title: "Tareas", Count: taskCount},
		{Title: "Materiales", Count: materialCount},
		{Title: "Productos", Count: productCount},
	}, nil
}
