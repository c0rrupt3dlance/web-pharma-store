package services

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
)

type Authorization interface {
	Create(user models.User) (int, error)
	GenerateToken(username, password string) (string, error)
	VerifyToken(tokenString string) (int, error)
}

type Service struct {
	Authorization
}

func NewService(repo repository.Repository) *Service {
	return &Service{
		New
	}
}
