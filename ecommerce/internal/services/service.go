package services

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type Products interface {
	Create(product models.Product) (int, error)
	GetById(ProductId int) (models.ProductResponse, error)
	Update(product models.Product) error
	Delete(ProductId int) error
}

type Service struct {
	Products
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Products: NewService(repo),
	}
}
