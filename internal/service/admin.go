package service

import (
	"fmt"

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
		return domain.ErrAdminSetupDisabled
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("setup admin: %w", err)
	}

	user := &domain.User{Username: username, Password: string(hashedPassword), Email: email, RoleID: 2}

	s.ur.Create(user)

	return nil
}
