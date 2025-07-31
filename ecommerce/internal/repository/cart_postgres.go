package repository

import (
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type CartPostgres struct {
	pool *pgxpool.Pool
}

func NewCartPostgres(pool *pgxpool.Pool) *CartPostgres {
	return &CartPostgres{
		pool,
	}
}

func (r *CartPostgres) AddToCart(p int, q int, price float32, userId int) error {
	query := fmt.Sprintf(`
		INSERT INTO cart_products (cart_id, product_id, quantity, price) 
		VALUES((SELECT id FROM user_carts where user_id=$1) $2,$3,$4) 
	`)

	return nil
}
