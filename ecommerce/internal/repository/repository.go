package repository

import (
	models "github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Products interface {
	Create(product models.Product) (int, error)
	GetById(ProductId int) (models.ProductResponse, error)
	Update(product models.Product) error
	Delete(ProductId int) error
}

type Repository struct {
	Products
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Products: NewProductPostgres(pool),
	}
}
