package helps

import "time"

type Help struct {
	ID          uint      `gorm:"primarykey" json:"id"`
	Tag         string    `gorm:"type:varchar(255);uniqueIndex;not null" json:"tag"`
	Description string    `gorm:"type:text;not null" json:"description"`
	Active      bool      `gorm:"default:true" json:"active"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Help) TableName() string {
	return "helps"
}
