package repository

import "context"

type Repository interface {
	GetOrders(ctx context.Context, parameter string) (models.Orders, error)
	CreateOrder(ctx context.Context, models.Order) error
}
