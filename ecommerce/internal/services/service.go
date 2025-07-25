package services

import (
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type Products interface {
	Create(ctx context.Context, product models.ProductInput) (int, error)
	GetById(ctx context.Context, ProductId int) (models.ProductResponse, error)
	Update(ctx context.Context, productId int, product models.UpdateProductInput) error
	Delete(ctx context.Context, ProductId int) error
	GetByCategories(ctx context.Context, categoriesId []int) ([]models.ProductResponse, error)
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

type FileStorage interface {
	Add(productId int, mediaFiles []models.FileDataType) (map[string]string, error)
}

type Service struct {
	Products
	Authorization
	Cart
	FileStorage
}

func NewService(repo *repository.Repository, signingKey string) *Service {
	return &Service{
		Products:      NewProductsService(repo),
		Authorization: NewAuthService(signingKey),
	}
}
