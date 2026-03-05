package support

import (
	"time"

	"github.com/TFX0019/api-go-gds/features/auth"
	"github.com/google/uuid"
)

type SupportCategory struct {
	ID          uuid.UUID `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string    `gorm:"type:text;not null"`
	Description string    `gorm:"type:text"`
	Active      bool      `gorm:"not null;default:true"`
	CreatedAt   time.Time `gorm:"not null;default:now()"`
	UpdatedAt   time.Time `gorm:"not null;default:now()"`
}

func (SupportCategory) TableName() string {
	return "support_categories"
}

type Support struct {
	ID                uuid.UUID       `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Subject           string          `gorm:"type:text;not null"`
	Description       string          `gorm:"type:text;not null"`
	UserID            uint            `gorm:"not null"`
	User              auth.User       `gorm:"foreignKey:UserID"`
	SupportCategoryID uuid.UUID       `gorm:"type:uuid;not null"`
	SupportCategory   SupportCategory `gorm:"foreignKey:SupportCategoryID"`
	Status            string          `gorm:"type:text;not null;default:'open'"`
	Image             string          `gorm:"type:text"`
	ParentID          *uuid.UUID      `gorm:"type:uuid"`
	Parent            *Support        `gorm:"foreignKey:ParentID"`
	IsDeleted         bool            `gorm:"not null;default:false"`
	CreatedAt         time.Time       `gorm:"not null;default:now()"`
	UpdatedAt         time.Time       `gorm:"not null;default:now()"`
	Replies           []Support       `gorm:"foreignKey:ParentID"`
}

func (Support) TableName() string {
	return "supports"
}
