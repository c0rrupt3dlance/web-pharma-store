package repository

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Authorization interface {
	Create(models.User) (int, error)
	GetUser(username, password string) (models.User, error)
}

type Repository struct {
	Authorization
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(pool),
	}
}
