package tasks

import (
	"time"

	"github.com/TFX0019/api-go-gds/features/products"
	"github.com/google/uuid"
)

type Task struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID      uint              `gorm:"not null"`
	Name        string            `gorm:"type:text;not null"`
	Description string            `gorm:"type:text"`
	Status      string            `gorm:"type:text;not null;check:status IN ('pending', 'in_progress', 'completed', 'canceled')"`
	DateTime    time.Time         `gorm:"type:timestamp with time zone;not null"`
	ProductID   *uuid.UUID        `gorm:"type:uuid"`
	Product     *products.Product `gorm:"foreignKey:ProductID"`
	CreatedAt   time.Time         `gorm:"not null;default:now()"`
	UpdatedAt   time.Time         `gorm:"not null;default:now()"`
}

func (Task) TableName() string {
	return "tasks"
}
