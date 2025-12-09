package products

import (
	"time"

	"github.com/google/uuid"
)

type Product struct {
	ID                   uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID               uint       `gorm:"not null"`
	Name                 string     `gorm:"type:text;not null"`
	ClientID             *uuid.UUID `gorm:"type:uuid"` // Optional
	MaterialsCost        float64    `gorm:"type:numeric;not null"`
	HoursCost            float64    `gorm:"type:numeric;not null"`
	ProfitPercentage     float64    `gorm:"type:numeric;not null"`
	IncludeFixedExpenses bool       `gorm:"type:boolean;not null"`
	FixedExpenseRate     float64    `gorm:"type:numeric;not null"`
	Subtotal             float64    `gorm:"type:numeric;not null"`
	FixedExpensesAmount  float64    `gorm:"type:numeric;not null"`
	BaseTotal            float64    `gorm:"type:numeric;not null"`
	ProfitAmount         float64    `gorm:"type:numeric;not null"`
	Total                float64    `gorm:"type:numeric;not null"`
	CreatedAt            time.Time  `gorm:"not null;default:now()"`
	UpdatedAt            time.Time  `gorm:"not null;default:now()"`
}

func (Product) TableName() string {
	return "products"
}
