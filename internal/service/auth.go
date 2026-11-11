package service

import (
	"crypto/rsa"
	"errors"
	"log/slog"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victorhslito/ticketsystem/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	ur        repository.UserRepository
	signKey   *rsa.PrivateKey
	verifyKey *rsa.PublicKey
	expiresIn time.Duration
	logger    *slog.Logger
}

type CustomClaims struct {
	Roles  []string `json:"roles"`
	Email  string   `json:"email"`
	UserId uint     `json:"userId"`
	jwt.RegisteredClaims
}

type AccessTokenResponseDTO struct {
	AccessToken string        `json:"accessToken"`
	ExpiresIn   time.Duration `json:"expiresIn"`
}

func NewAuthService(
	userRepository repository.UserRepository,
	signKey *rsa.PrivateKey,
	verifyKey *rsa.PublicKey,
	expiresIn time.Duration,
	logger *slog.Logger,
) *AuthService {
	return &AuthService{
		ur:        userRepository,
		signKey:   signKey,
		verifyKey: verifyKey,
		expiresIn: expiresIn,
		logger:    logger.With("component", "auth_service"),
	}
}

func (s *AuthService) GenerateToken(username string, password string) (*AccessTokenResponseDTO, error) {
	user, err := s.ur.FindByUsername(username)
	if err != nil {
		s.logger.Info("User not found")
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		s.logger.Info("Passwords doesn't match")
		return nil, errors.New("password doesn't match")
	}

	issuedAt := jwt.NewNumericDate(time.Now())

	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Name))

	t.Claims = &CustomClaims{
		[]string{"User"},
		user.Email,
		user.ID,
		jwt.RegisteredClaims{
			Issuer:    "ticketsystem",
			Subject:   username,
			IssuedAt:  issuedAt,
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(s.expiresIn * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, t.Claims)
	ss, err := token.SignedString(s.signKey)

	return &AccessTokenResponseDTO{
		AccessToken: ss,
		ExpiresIn:   s.expiresIn,
	}, nil
}
