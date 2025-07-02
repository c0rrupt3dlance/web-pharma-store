package services

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
)

type Authorization interface {
	Create(user models.User) (int, error)
	GenerateTokens(username, password string) (string, string, error)
	RefreshTokens(refreshToken string) (string, string, error)
	VerifyAccessToken(tokenString string) (models.User, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		NewAuthService(repo),
	}
}
