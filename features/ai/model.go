package ai

import (
	"time"

	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/google/uuid"
)

type AIGeneration struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	UserID      uint      `gorm:"not null" json:"user_id"`
	User        auth.User `gorm:"foreignKey:UserID" json:"-"`
	Prompt      string    `gorm:"type:text;not null" json:"prompt"`
	ImageInput  *string   `gorm:"type:text" json:"image_input"`
	ImageOutput *string   `gorm:"type:text" json:"image_output"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (AIGeneration) TableName() string {
	return "ai_generations"
}
