package repository

import (
	"context"
	models "github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Products interface {
	Create(ctx context.Context, product models.ProductInput) (int, error)
	GetById(ctx context.Context, productId int) (models.ProductResponse, error)
	Update(ctx context.Context, productId int, product models.UpdateProductInput) error
	Delete(ctx context.Context, ProductId int) error
	GetByCategories(ctx context.Context, categoriesId []int) ([]models.ProductResponse, error)
}

type Cart interface {
	AddItem(userId, productId, quantity int) error
	UpdateQuantity(userId, productId, quantity int) error
	RemoveItem(userId, productId int) error
	GetCart(userId int) ([]models.CartItem, error)
	ClearCart(userId int) error
}
type Repository struct {
	Products
	Cart
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Products: NewProductPostgres(pool),
	}
}
