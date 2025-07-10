package repository

import (
	"context"
	"fmt"
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
	query := fmt.Sprintf("INSERT INTO products (name, description, price) values ($1, $2, $3) returning id")
	row := r.pool.QueryRow(context.Background(), query, product.Name, product.Description, product.Price)
	if err := row.Scan(&product.Id); err != nil {
		return 0, err
	}
	return product.Id, nil
}
func (r *ProductPostgres) GetById(ProductId int) (models.ProductResponse, error) {
	var product models.ProductResponse
	query := fmt.Sprintf(`select * from products where id=$1`)
	row := r.pool.QueryRow(context.Background(), query, ProductId)
	if err := row.Scan(&product.Id, &product.Name, &product.Description, &product.Price); err != nil {
		return models.ProductResponse{}, err
	}
	return product, nil
}
func (r *ProductPostgres) Update(product models.Product) error {
	return nil
}
func (r *ProductPostgres) Delete(ProductId int) error {
	return nil
}
