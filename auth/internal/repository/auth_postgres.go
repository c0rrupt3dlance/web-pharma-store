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
	query := fmt.Sprintf("insert into %s (name, surname, email, password) values ($1, $2, $3, $4) returning id", usersTable)

	row := r.pool.QueryRow(context.Background(), query, user.Name, user.Surname, user.Email, user.Password)
	var userId int

	if err := row.Scan(&userId); err != nil {
		return 0, err
	}

	return userId, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	query := fmt.Sprintf("select id, name, surname, email, is_admin from %s where username = $1 and password = $2", usersTable)

	row := r.pool.QueryRow(context.Background(), query, username, password)
	var user models.User
	user.Username = username
	if err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Email, &user.IsAdmin); err != nil {
		return models.User{}, err
	}

	return user, nil
}
