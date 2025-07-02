package repository

import (
	"context"
	"fmt"
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
	query := fmt.Sprintf("INSERT INTO %s (name, surname, username, email, password) values ($1,$2,$3,$4,$5) returning id", usersTable)
	err := r.pool.QueryRow(context.Background(), query, user.Name, user.Surname, user.Username, user.Email, user.Password).Scan(&user.Id)
	if err != nil {
		return 0, err
	}
	return user.Id, nil
}

func (r *AuthPostgres) GetUser(username string) (models.User, error) {
	user := models.User{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE username = $1", usersTable)

	row := r.pool.QueryRow(context.Background(), query, username)

	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Username, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *AuthPostgres) SaveRefreshToken(token models.RefreshToken) error {
	query := fmt.Sprintf(`insert into %s (user_id, token, expires_at, revoked)
		values ($1, $2, $3, $4)`, refreshTokensTable)

	_, err := r.pool.Exec(context.Background(), query,
		token.UserId, token.Token, token.ExpiresAt, token.Revoked)
	if err != nil {
		return err
	}

	return nil
}

func (r *AuthPostgres) GetRefreshToken(token string) (models.RefreshToken, error) {
	refreshToken := models.RefreshToken{}
	query := fmt.Sprintf(`select user_id, expires_at, revoked, issued_at 
		from %s where token_string=$1`, refreshTokensTable)

	err := r.pool.QueryRow(context.Background(), query, token).Scan(&refreshToken.UserId,
		&refreshToken.ExpiresAt, &refreshToken.Revoked, &refreshToken.IssuedAt)
	if err != nil {
		return models.RefreshToken{}, err
	}

	return refreshToken, nil
}

func (r *AuthPostgres) GetUserById(userId int) (models.User, error) {
	user := models.User{}
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", usersTable)

	row := r.pool.QueryRow(context.Background(), query, userId)

	err := row.Scan(&user.Id, &user.Name, &user.Surname, &user.Username, &user.Email, &user.Password, &user.Role)

	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r *AuthPostgres) RevokeRefreshToken(tokenString string) error {
	query := fmt.Sprintf(`insert into %s set revoked=true where token=$2`, refreshTokensTable)

	_, err := r.pool.Exec(context.Background(), query, tokenString)
	if err != nil {
		return err
	}

	return nil
}
