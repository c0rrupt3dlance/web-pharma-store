package repository

import (
	"context"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type CartPostgres struct {
	pool *pgxpool.Pool
}

func NewCartPostgres(pool *pgxpool.Pool) *CartPostgres {
	return &CartPostgres{
		pool,
	}
}

func (r *CartPostgres) GetUserCart(ctx context.Context, userId int) (int, error) {
	var cartId int
	query := fmt.Sprintf(`
	SELECT id from %s where user_id
	`, userCartsTable)
	row := r.pool.QueryRow(ctx, query, userId)
	if err := row.Scan(&cartId); err != nil {
		logrus.Println("GetUserCart:", err)
		return 0, err
	}

	return 0, nil

}

func (r *CartPostgres) AddToCart(p int, q int, price float32, userId int) error {
	query := fmt.Sprintf(`
		INSERT INTO cart_products (cart_id, product_id, quantity, price) 
		VALUES((SELECT id FROM user_carts where user_id=$1) $2,$3,$4) 
	`)

	return nil
}
