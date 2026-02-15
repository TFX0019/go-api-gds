package dashboard

import (
	"github.com/TFX0019/api-go-gds/features/customers"
	"github.com/TFX0019/api-go-gds/features/materials"
	"github.com/TFX0019/api-go-gds/features/products"
	"github.com/TFX0019/api-go-gds/features/tasks"
	"gorm.io/gorm"
)

type Repository interface {
	GetSummaryCounts() ([]SummaryItem, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) GetSummaryCounts() ([]SummaryItem, error) {
	var customerCount int64
	var taskCount int64
	var materialCount int64
	var productCount int64

	// Count all customers
	if err := r.db.Model(&customers.Customer{}).Count(&customerCount).Error; err != nil {
		return nil, err
	}

	// Count all tasks
	if err := r.db.Model(&tasks.Task{}).Count(&taskCount).Error; err != nil {
		return nil, err
	}

	// Count all materials
	if err := r.db.Model(&materials.Material{}).Count(&materialCount).Error; err != nil {
		return nil, err
	}

	// Count all products
	if err := r.db.Model(&products.Product{}).Count(&productCount).Error; err != nil {
		return nil, err
	}

	return []SummaryItem{
		{Title: "Clientes", Count: customerCount},
		{Title: "Tareas", Count: taskCount},
		{Title: "Materiales", Count: materialCount},
		{Title: "Productos", Count: productCount},
	}, nil
}
