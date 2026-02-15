package plans

import (
	"time"
)

type Plan struct {
	ID           uint      `gorm:"primarykey" json:"id"`
	ProductID    string    `gorm:"uniqueIndex;not null" json:"product_id"`
	Title        string    `gorm:"not null" json:"title"`
	Description  string    `gorm:"not null" json:"description"`
	Price        float64   `gorm:"type:numeric;not null" json:"price"`
	Benefits     []string  `gorm:"serializer:json" json:"benefits"`
	MaxCustomers int       `gorm:"not null;default:20" json:"max_customers"` // -1 for unlimited
	MaxProducts  int       `gorm:"not null;default:20" json:"max_products"`  // -1 for unlimited
	MaxMaterials int       `gorm:"not null;default:20" json:"max_materials"` // -1 for unlimited
	MaxTasks     int       `gorm:"not null;default:20" json:"max_tasks"`     // -1 for unlimited
	IsActive     bool      `gorm:"default:true" json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}
