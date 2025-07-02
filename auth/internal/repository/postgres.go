package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	usersTable         = "users"
	refreshTokensTable = "refresh_tokens"
)

type PgPool struct {
	pool *pgxpool.Pool
}

type PgConfig struct {
	Host     string
	Port     string
	Database string
	User     string
	Password string
}

func NewPgPool(cfg PgConfig) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		return nil, errors.New("unable to connect to postgres db")
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, err
	}

	return pool, err
}
