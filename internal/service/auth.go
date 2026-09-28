package service

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/victorhslito/ticketsystem/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	expiresIn      = 300
	privateKeyPath = "../rsaPrivateKey.pem"
	publicKeyPath  = "../publicKey.pem"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type AuthService struct {
	ur repository.UserRepository
}

type CustomClaims struct {
	Roles []string `json:"roles"`
	Email string   `json:"email"`
	jwt.RegisteredClaims
}

type AccessTokenResponseDTO struct {
	AccessToken string        `json:"accessToken"`
	ExpiresIn   time.Duration `json:"expiresIn"`
}

func NewAuthService(userRepository repository.UserRepository) *AuthService {
	initJwtConfig()
	return &AuthService{ur: userRepository}
}

func initJwtConfig() {
	signBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Fatal("Erro ao ler chave privada: ", err)
	}

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	if err != nil {
		log.Fatal("Erro ao fazer parse da chave privada: ", err)
	}

	verifyBytes, err := os.ReadFile(publicKeyPath)
	if err != nil {
		log.Fatal("Erro ao ler chave pública: ", err)
	}

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	if err != nil {
		log.Fatal("Erro ao fazer parse da chave pública: ", err)
	}
}

func (s *AuthService) GenerateToken(username string, password string) (*AccessTokenResponseDTO, error) {
	user, err := s.ur.FindByUsername(username)
	if err != nil {
		fmt.Printf("User not found!\n")
		return nil, errors.New("user not found")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		fmt.Printf("Passwords doesn't match\n")
		return nil, errors.New("password doesn't match")
	}

	issuedAt := jwt.NewNumericDate(time.Now())

	t := jwt.New(jwt.GetSigningMethod(jwt.SigningMethodRS256.Name))

	t.Claims = &CustomClaims{
		[]string{"User", "Admin"},
		user.Email,
		jwt.RegisteredClaims{
			Issuer:    "ticketsystem",
			Subject:   username,
			IssuedAt:  issuedAt,
			ExpiresAt: jwt.NewNumericDate(issuedAt.Add(expiresIn * time.Second)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, t.Claims)
	ss, err := token.SignedString(signKey)

	return &AccessTokenResponseDTO{
		AccessToken: ss,
		ExpiresIn:   expiresIn,
	}, nil
}
