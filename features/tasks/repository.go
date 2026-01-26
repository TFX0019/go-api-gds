package tasks

import "gorm.io/gorm"

type Repository interface {
	Create(task *Task) error
	FindAll(limit, offset int, status, date string) ([]Task, int64, error)
	FindByID(id string) (*Task, error)
	FindByUserID(userID string, limit, offset int, status, date string) ([]Task, int64, error)
	Update(task *Task) error
	Delete(id string) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(task *Task) error {
	// Preload to return full object if needed, but usually Create just saves.
	// We might need to fetch it back to return full structure with preloads if requested.
	// For now, simple create.
	return r.db.Create(task).Error
}

func (r *repository) FindAll(limit, offset int, status, date string) ([]Task, int64, error) {
	var tasks []Task
	var total int64

	db := r.db.Model(&Task{})

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if date != "" {
		// Asumming date_time is timestamp. Filter by day.
		// date format YYYY-MM-DD
		db = db.Where("DATE(date_time) = ?", date)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Preload("Product.Client").Limit(limit).Offset(offset).Order("date_time desc").Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *repository) FindByID(id string) (*Task, error) {
	var task Task
	err := r.db.Preload("Product.Client").Where("id = ?", id).First(&task).Error
	if err != nil {
		return nil, err
	}
	return &task, nil
}

func (r *repository) FindByUserID(userID string, limit, offset int, status, date string) ([]Task, int64, error) {
	var tasks []Task
	var total int64

	db := r.db.Model(&Task{}).Where("user_id = ?", userID)

	if status != "" {
		db = db.Where("status = ?", status)
	}
	if date != "" {
		db = db.Where("DATE(date_time) = ?", date)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Preload("Product.Client").Limit(limit).Offset(offset).Order("date_time desc").Find(&tasks).Error
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *repository) Update(task *Task) error {
	return r.db.Save(task).Error
}

func (r *repository) Delete(id string) error {
	return r.db.Delete(&Task{}, "id = ?", id).Error
}
