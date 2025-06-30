package repository

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type AuthPostgres struct {
	pool *pgxpool.Pool
}

func NewAuthPostgres(pool *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{pool: pool}
}

func (r *AuthPostgres) Create(user models.User) (int, error) {
	return 0, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	return models.User{}, nil
}
