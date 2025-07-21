package repository

import "github.com/jackc/pgx/v5/pgxpool"

type CartPostgres struct {
	pool *pgxpool.Pool
}

func NewCartPostgres(pool *pgxpool.Pool) *CartPostgres {
	return &CartPostgres{
		pool,
	}
}

func (r *CartPostgres) AddToCart() error {
	return nil
}
