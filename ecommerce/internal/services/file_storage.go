package services

import (
	"context"
	"errors"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FileStorageService struct {
	media    repository.FileStorage
	products repository.Products
}

func NewFileStorageService(media repository.FileStorage, products repository.Products) *FileStorageService {
	return &FileStorageService{
		media:    media,
		products: products,
	}
}

func (s *FileStorageService) AddMedia(ctx context.Context, productId int, mediaFiles []models.FileDataType) ([]models.MediaUrl, error) {
	objIds := make(map[int]string)
	data := make(map[string]models.FileDataType)
	for _, v := range mediaFiles {
		objectId := uuid.New().String()
		objIds[v.Position] = objectId
		data[objectId] = v
	}

	err := s.products.AddProductMedia(ctx, productId, objIds)
	if err != nil {
		logrus.Println("unable to add to the products_media table, reason:", err)
		return nil, errors.New("unable to add media")
	}

	urls, err := s.media.AddMedia(data)
	if err != nil {
		return nil, err
	}
	var responseUrls []models.MediaUrl
	for k, v := range objIds {
		responseUrls = append(responseUrls, models.MediaUrl{
			Url:      urls[v],
			Position: k,
		})
	}

	return responseUrls, nil
}

func (s *FileStorageService) GetMedia(ctx context.Context, productId int) ([]models.MediaUrl, error) {
	productMedia := make([]models.MediaUrl, 0)

	objectIds, err := s.products.GetProductMedia(ctx, productId)
	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	urls, err := s.media.GetMedia(ctx, objectIds)
	if err != nil {
		logrus.Println(err)
		return nil, fmt.Errorf("Error while %w:", err)
	}

	for _, v := range objectIds {
		var media = models.MediaUrl{
			Url: urls[v],
		}
		productMedia = append(productMedia, media)
	}

	return productMedia, nil
}
