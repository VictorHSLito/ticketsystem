package service

import (
	"errors"
	"log"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	r repository.UserRepository
}

func NewUserService(r repository.UserRepository) *UserService {
	return &UserService{r: r}
}

func (s *UserService) CreateUser(username, email, password string) (uint, error) {
	if s.isUserAlreadyExists(email) {
		return 0, errors.New("user already exists")
	}
	user := &domain.User{
		Username: username,
		Email:    email,
		RoleID:   1,
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("error trying to encrypt password, err: ", err.Error())
	}

	user.Password = string(hashedPassword)

	s.r.Create(user)

	return user.ID, nil
}

func (s *UserService) FindByID(userID uint) *domain.User {
	user, err := s.r.ListByID(userID)
	if err != nil {
		return nil
	}

	return user
}

func (s *UserService) isUserAlreadyExists(email string) bool {
	user, err := s.r.ListByEmail(email)

	if user != nil {
		return true
	}

	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) {
		return false
	}

	return false
}
