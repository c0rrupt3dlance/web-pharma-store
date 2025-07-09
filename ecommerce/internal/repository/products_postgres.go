package repository

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductPostgres struct {
	pool *pgxpool.Pool
}

func NewProductPostgres(pool *pgxpool.Pool) *ProductPostgres {
	return &ProductPostgres{
		pool: pool,
	}
}

func (r *ProductPostgres) Create(product models.Product) (int, error) {
	return 0, nil
}
func (r *ProductPostgres) GetById(ProductId int) (models.ProductResponse, error) {
	return models.ProductResponse{}, nil
}
func (r *ProductPostgres) Update(product models.Product) error {
	return nil
}
func (r *ProductPostgres) Delete(ProductId int) error {
	return nil
}
