package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"time"
)

type Event struct {
	Id            int       `json:"id"`
	AggregateType string    `json:"aggregate_type"`
	AggregateId   int       `json:"aggregate_id"`
	EventType     string    `json:"event_type"`
	Payload       []byte    `json:"payload"`
	CreatedAt     time.Time `json:"created_at"`
}

const (
	outboxTable = "outbox"
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

func NewPgPool(cfg PgConfig, ctx context.Context) (*pgxpool.Pool, error) {
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database))
	if err != nil {
		logrus.Printf("fail due to %s", err)
		return nil, errors.New("unable to connect to postgres db")
	}

	err = pool.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return pool, err
}

type Repository struct {
	pool *pgxpool.Pool
}

func NewMQPostgres(pool *pgxpool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) FetchUnprocessed(ctx context.Context, limit int) ([]Event, error) {
	query := fmt.Sprintf(`
	SELECT id, aggregate_type, aggregate_id, 
	       event_type, payload, created_at
	FROM %s
	WHERE processed_at IS NULL
	ORDER BY created_at
	LIMIT $1
	`, outboxTable)
	rows, err := r.pool.Query(ctx, query, limit)
	if err != nil {
		return nil, errors.New("couldn't get events from outbox")
	}
	defer rows.Close()

	var events []Event
	for rows.Next() {
		var e Event
		if err = rows.Scan(&e.Id, &e.AggregateType, &e.AggregateId,
			&e.EventType, &e.Payload, &e.CreatedAt); err != nil {
			return nil, errors.New("couldn't scan row")
		}
		events = append(events, e)
	}

	return events, nil
}

func (r *Repository) MarkProcessed(ctx context.Context, id int) error {
	query := fmt.Sprintf(`
	UPDATE %s 
	SET processed_at = now() 
	WHERE id = $1
	`, outboxTable)

	_, err := r.pool.Exec(ctx, query, id)
	return err
}
