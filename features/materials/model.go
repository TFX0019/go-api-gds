package materials

import (
	"time"

	"github.com/google/uuid"
)

type Material struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID    uint      `gorm:"not null"`
	Name      string    `gorm:"type:text;not null"`
	Price     float64   `gorm:"type:numeric;not null"`
	Unit      string    `gorm:"type:text;not null"`
	ImageURL  string    `gorm:"type:text"`
	CreatedAt time.Time `gorm:"not null;default:now()"`
	UpdatedAt time.Time `gorm:"not null;default:now()"`
}

func (Material) TableName() string {
	return "materials"
}
