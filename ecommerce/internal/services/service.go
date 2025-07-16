package services

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type Products interface {
	Create(product models.ProductInput) (int, error)
	GetById(ProductId int) (models.ProductResponse, error)
	Update(product models.UpdateProductInput) error
	Delete(ProductId int) error
}

type Authorization interface {
	VerifyAccessToken(accessToken string) (int, error)
}

type Cart interface {
	AddItem(userId, productId, quantity int) error
	UpdateQuantity(userId, productId, quantity int) error
	RemoveItem(userId, productId int) error
	GetCart(userId int) ([]models.CartItem, error)
	ClearCart(userId int) error
}

type Service struct {
	Products
	Authorization
	Cart
}

func NewService(repo *repository.Repository, signingKey string) *Service {
	return &Service{
		Products:      NewProductsService(repo),
		Authorization: NewAuthService(signingKey),
	}
}
