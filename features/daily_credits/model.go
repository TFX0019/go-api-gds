package daily_credits

import "time"

type DailyCredit struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Free      int       `gorm:"not null;default:3" json:"free"`
	Premium   int       `gorm:"not null;default:6" json:"premium"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (DailyCredit) TableName() string {
	return "daily_credits"
}
