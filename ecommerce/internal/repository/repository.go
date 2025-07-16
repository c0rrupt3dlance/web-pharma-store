package repository

import (
	models "github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Products interface {
	Create(product models.ProductInput) (int, error)
	GetById(productId int) (models.ProductResponse, error)
	Update(product models.UpdateProductInput) error
	Delete(ProductId int) error
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
