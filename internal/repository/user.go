package repository

import (
	"errors"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"gorm.io/gorm"
)

type UserRepository interface {
	ListByID(userID uint) (*domain.User, error)
	ListByEmail(email string) (*domain.User, error)
	Create(user *domain.User) (*domain.User, error)
	Update(user *domain.User) (*domain.User, error)
}

type userRepositoryImpl struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{db: db}
}

func (r *userRepositoryImpl) ListByID(userID uint) (*domain.User, error) {
	var user domain.User

	err := r.db.Where("id = ?", userID).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) ListByEmail(email string) (*domain.User, error) {
	var user domain.User

	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	return &user, nil
}

func (r *userRepositoryImpl) Create(user *domain.User) (*domain.User, error) {
	err := r.db.Create(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *userRepositoryImpl) Update(user *domain.User) (*domain.User, error) {
	err := r.db.Save(user).Error
	if err != nil {
		return nil, err
	}

	return user, nil
}
