package products

import (
	"time"

	"github.com/TFX0019/api-go-gds/features/customers"
	"github.com/google/uuid"
)

type Product struct {
	ID                   uuid.UUID           `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               uint                `gorm:"not null"`
	Name                 string              `gorm:"type:text;not null"`
	ClientID             *uuid.UUID          `gorm:"type:uuid"` // Optional
	Client               *customers.Customer `gorm:"foreignKey:ClientID"`
	MaterialsCost        float64             `gorm:"type:numeric;not null"`
	HoursCost            float64             `gorm:"type:numeric;not null"`
	ProfitPercentage     float64             `gorm:"type:numeric;not null"`
	IncludeFixedExpenses bool                `gorm:"type:boolean;not null"`
	FixedExpenseRate     float64             `gorm:"type:numeric;not null"`
	Subtotal             float64             `gorm:"type:numeric;not null"`
	FixedExpensesAmount  float64             `gorm:"type:numeric;not null"`
	BaseTotal            float64             `gorm:"type:numeric;not null"`
	ProfitAmount         float64             `gorm:"type:numeric;not null"`
	Total                float64             `gorm:"type:numeric;not null"`
	Status               string              `gorm:"type:text;not null;default:'pending'"`
	DatePaid             *time.Time          `gorm:"type:timestamp"`
	Images               []ProductImage      `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	CreatedAt            time.Time           `gorm:"not null;default:now()"`
	UpdatedAt            time.Time           `gorm:"not null;default:now()"`
}

func (Product) TableName() string {
	return "products"
}

type ProductImage struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	ProductID uuid.UUID `gorm:"type:uuid;not null"`
	Path      string    `gorm:"type:text;not null"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
}

func (ProductImage) TableName() string {
	return "product_images"
}
