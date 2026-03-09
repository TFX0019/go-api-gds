package banners

import (
	"time"

	"github.com/google/uuid"
)

type Banner struct {
	ID        uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Image     string    `gorm:"type:text;not null" json:"image"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Banner) TableName() string {
	return "banners"
}
