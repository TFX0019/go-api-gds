package daily_credits

import (
	"gorm.io/gorm"
)

type Repository interface {
	Get() (*DailyCredit, error)
	Update(credit *DailyCredit) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Get() (*DailyCredit, error) {
	var config DailyCredit
	err := r.db.First(&config).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Create a default one if not exists
			defaultConfig := DailyCredit{Free: 3, Premium: 6}
			if createErr := r.db.Create(&defaultConfig).Error; createErr != nil {
				return nil, createErr
			}
			return &defaultConfig, nil
		}
		return nil, err
	}
	return &config, nil
}

func (r *repository) Update(credit *DailyCredit) error {
	return r.db.Save(credit).Error
}
