package services

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
)

type ProductsService struct {
	repo repository.Products
}

func NewProductsService(repo *repository.Repository) *ProductsService {
	return &ProductsService{
		repo: repo,
	}
}

func (s *ProductsService) Create(product models.ProductInput) (int, error) {
	return s.repo.Create(product)
}
func (s *ProductsService) GetById(ProductId int) (models.ProductResponse, error) {
	return s.repo.GetById(ProductId)
}
func (s *ProductsService) Update(product models.UpdateProductInput) error {
	return s.repo.Update(product)
}
func (s *ProductsService) Delete(ProductId int) error {
	return s.repo.Delete(ProductId)
}
