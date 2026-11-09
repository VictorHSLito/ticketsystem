package service

import (
	"fmt"
	"log"

	"github.com/victorhslito/ticketsystem/internal/domain"
	"github.com/victorhslito/ticketsystem/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
	ur repository.UserRepository
}

func NewAdminService(ur repository.UserRepository) *AdminService {
	return &AdminService{ur: ur}
}

func (s *AdminService) SetupAdminUser(username, password, email string) error {
	users := s.ur.CountUsers()

	if users > 0 {
		return fmt.Errorf("Admin user can only be created where there are no users in the system.")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Could not be possible to encrypt password")
	}

	user := &domain.User{Username: username, Password: string(hashedPassword), Email: email, RoleID: 2}

	s.ur.Create(user)

	return nil
}
