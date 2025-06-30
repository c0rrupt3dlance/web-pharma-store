package services

import (
	"crypto/sha1"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"time"
)

const (
	salt       = "M!m39233j4!"
	signingKey = "(!#Mn45ntg24gN!0"
	TokenTTL   = time.Hour * 8
)

type tokenClaims struct {
	userId int
	jwt.RegisteredClaims
}

type AuthService struct {
	repo repository.Repository
}

func NewAuthService(repo repository.Repository) *AuthService {
	return &AuthService{repo: repo}
}

func generatePasswordHash(password string) string {
	passwordHash := sha1.New()
	passwordHash.Write([]byte(password))
	passwordHash.Write([]byte(salt))
	return fmt.Sprintf("%x", passwordHash.Sum(nil))
}

func (s *AuthService) Create(user models.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.Create(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, password)

	if err != nil {
		return "", err
	}

	claims := tokenClaims{
		user.Id,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "Pharma-Auth-Service",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(signingKey))

	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) VerifyToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(signingKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(tokenClaims)
	if !ok || !token.Valid {
		return 0, err
	}

	return claims.userId, nil
}
