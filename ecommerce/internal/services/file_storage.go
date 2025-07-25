package services

import "github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"

type FileStorageService struct {
	m repository.FileStorage
	r repository.Products
}
