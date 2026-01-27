package customers

import (
	"time"

	"github.com/google/uuid"
)

type Customer struct {
	ID               uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	UserID           uint      `gorm:"not null"`
	AvatarURL        string    `gorm:"type:text"`
	Name             string    `gorm:"type:text;not null"`
	Phone            string    `gorm:"type:text"`
	Email            string    `gorm:"type:text"`
	UsesStandardSize bool      `gorm:"type:boolean;not null"`
	StandardSize     string    `gorm:"type:text"`
	Back             *float64  `gorm:"type:numeric"`
	Neck             *float64  `gorm:"type:numeric"`
	FrontSize        *float64  `gorm:"type:numeric"`
	Armhole          *float64  `gorm:"type:numeric"`
	BackSize         *float64  `gorm:"type:numeric"`
	BustChest        *float64  `gorm:"type:numeric"`
	Waist            *float64  `gorm:"type:numeric"`
	Hip              *float64  `gorm:"type:numeric"`
	RiseHeight       *float64  `gorm:"type:numeric"`
	SkirtLength      *float64  `gorm:"type:numeric"`
	PantsLength      *float64  `gorm:"type:numeric"`
	KneeWidth        *float64  `gorm:"type:numeric"`
	HemWidth         *float64  `gorm:"type:numeric"`
	SleeveLength     *float64  `gorm:"type:numeric"`
	CuffSize         *float64  `gorm:"type:numeric"`
	CreatedAt        time.Time `gorm:"not null;default:now()"`
	UpdatedAt        time.Time `gorm:"not null;default:now()"`
}

func (Customer) TableName() string {
	return "customers"
}
