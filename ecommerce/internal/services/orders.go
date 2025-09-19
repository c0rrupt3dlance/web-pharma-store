package services

import (
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type OrderService struct {
	repo repository.Orders
}

func NewOrderService(repo repository.Orders) *OrderService {
	return &OrderService{
		repo: repo,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, input models.OrderInput) (int, error) {
	return s.repo.CreateOrder(ctx, input)
}

func (s *OrderService) GetOrder(ctx context.Context, clientId int) (models.Order, error) {
	return s.repo.GetOrder(ctx, clientId)
}

func (s *OrderService) GetAllOrders(ctx context.Context, clientId int) ([]models.Order, error) {
	return nil, nil
}
