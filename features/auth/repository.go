package auth

import (
	"gorm.io/gorm"
)

type Repository interface {
	CreateUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uint) (*User, error)
	FindVerifyToken(token string) (*User, error)
	UpdateUser(user *User) error
	CreateVerificationCode(code *VerificationCode) error
	FindVerificationCode(email, code string) (*VerificationCode, error)
	DeleteVerificationCode(email string) error
	CreateSession(session *Session) error
	FindSessionByToken(token string) (*Session, error)
	RevokeSession(token string) error
	FindRoleByName(name string) (*Role, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) CreateUser(user *User) error {
	return r.db.Create(user).Error
}

func (r *repository) FindByEmail(email string) (*User, error) {
	var user User
	err := r.db.Preload("Subscription").Preload("Roles").Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindByID(id uint) (*User, error) {
	var user User
	err := r.db.Preload("Subscription").Preload("Roles").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindVerifyToken(token string) (*User, error) {
	var user User
	err := r.db.Where("verification_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) UpdateUser(user *User) error {
	return r.db.Save(user).Error
}

func (r *repository) CreateVerificationCode(code *VerificationCode) error {
	return r.db.Create(code).Error
}

func (r *repository) FindVerificationCode(email, code string) (*VerificationCode, error) {
	var vc VerificationCode
	err := r.db.Where("email = ? AND code = ?", email, code).First(&vc).Error
	if err != nil {
		return nil, err
	}
	return &vc, nil
}

func (r *repository) DeleteVerificationCode(email string) error {
	return r.db.Where("email = ?", email).Delete(&VerificationCode{}).Error
}

func (r *repository) CreateSession(session *Session) error {
	return r.db.Create(session).Error
}

func (r *repository) FindSessionByToken(token string) (*Session, error) {
	var session Session
	err := r.db.Where("token = ? AND is_valid = ?", token, true).First(&session).Error
	if err != nil {
		return nil, err
	}
	return &session, nil
}

func (r *repository) RevokeSession(token string) error {
	return r.db.Model(&Session{}).Where("token = ?", token).Update("is_valid", false).Error
}

func (r *repository) FindRoleByName(name string) (*Role, error) {
	var role Role
	err := r.db.Where("name = ?", name).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}
