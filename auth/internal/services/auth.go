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
	TokenTTL   = time.Minute * 60
)

type tokenClaims struct {
	userId   int
	username string
	role     string
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

func (s *AuthService) GenerateAccessToken(username, password string) (string, error) {
	passwordHash, err := generatePasswordHash(password)
	if err != nil {
		return "", err
	}
	user, err := s.repo.GetUser(username, passwordHash)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		userId:   user.Id,
		username: user.Username,
		role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(signingKey))
}

func (s *AuthService) VerifyAccessToken(tokenString string) (models.User, error) {
	user := models.User{}
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(signingKey), nil
	})

	if err != nil {
		return models.User{}, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		user.Id = claims.userId
		user.Username = claims.username
		user.Role = claims.role
		return user, nil
	}

	return models.User{}, err
}
