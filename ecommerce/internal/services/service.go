package services

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type Products interface {
	Create(product models.ProductInput) (int, error)
	GetById(ProductId int) (models.ProductResponse, error)
	Update(product models.Product) error
	Delete(ProductId int) error
}

type Authorization interface {
	VerifyAccessToken(accessToken string) (int, error)
}

type Service struct {
	Products
	Authorization
}

func NewService(repo *repository.Repository, signingKey string) *Service {
	return &Service{
		Products:      NewProductsService(repo),
		Authorization: NewAuthService(signingKey),
	}
}
