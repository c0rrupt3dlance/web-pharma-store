package services

import (
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type ProductsService struct {
	repo repository.Products
}

func NewProductsService(repo repository.Products) *ProductsService {
	return &ProductsService{
		repo: repo,
	}
}

func (s *ProductsService) Create(ctx context.Context, product models.ProductInput) (int, error) {
	return s.repo.Create(ctx, product)
}
func (s *ProductsService) GetById(ctx context.Context, ProductId int) (models.ProductResponse, error) {
	return s.repo.GetById(ctx, ProductId)
}
func (s *ProductsService) Update(ctx context.Context, productId int, product models.UpdateProductInput) error {
	return s.repo.Update(ctx, productId, product)
}
func (s *ProductsService) Delete(ctx context.Context, ProductId int) error {
	return s.repo.Delete(ctx, ProductId)
}

func (s *ProductsService) GetByCategories(ctx context.Context, categoriesId []int) ([]models.ProductResponse, error) {
	return s.repo.GetByCategories(ctx, categoriesId)
}
