package repository

import (
	"context"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable = "users"
)

type AuthPostgres struct {
	pool *pgxpool.Pool
}

func NewAuthPostgres(pool *pgxpool.Pool) *AuthPostgres {
	return &AuthPostgres{pool: pool}
}

func (r *AuthPostgres) Create(user models.User) (int, error) {
	query := fmt.Sprintf("INSERT INTO %s (name, surname, username, password, role) values ($1,$2,$3,$4,$5) returning id", usersTable)
	err := r.pool.QueryRow(context.Background(), query, user.Name, user.Surname, user.Username, user.Role).Scan(&user.Id)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	user := models.User{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1 AND password = $2", usersTable)

	row := r.pool.QueryRow(context.Background(), query, username, password)

	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Role)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}
