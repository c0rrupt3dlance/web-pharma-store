package services

import (
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	salt       = "M!m39233j4!"
	signingKey = "(!#Mn45ntg24gN!0"
	TokenTTL   = time.Hour * 8
)

type tokenClaims struct {
	userId   int
	username string
	isAdmin  bool
	jwt.RegisteredClaims
}

type AuthService struct {
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func generatePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func (s *AuthService) Create(user models.User) (int, error) {
	var err error
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	passwordHash, err := generatePasswordHash(password)
	user, err := s.repo.GetUser(username, passwordHash)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		userId:   user.Id,
		username: user.Username,
		isAdmin:  user.IsAdmin,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) ValidateToken(tokenString string) (models.User, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(signingKey), nil
	})
	if err != nil {
		return models.User{}, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return models.User{
			Id:       int(claims["userId"].(int)),
			Username: claims["username"].(string),
			IsAdmin:  claims["isAdmin"].(bool),
		}, nil
	}

	return models.User{}, err
}
