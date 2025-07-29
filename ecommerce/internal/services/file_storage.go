package services

import (
	"context"
	"errors"
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
	keys := make([]string, 0)
	data := make(map[string]models.FileDataType)
	for _, v := range mediaFiles {
		objectId := uuid.New().String()
		keys = append(keys, objectId)
		data[objectId] = v
	}

	err := s.products.AddProductMedia(ctx, productId, keys)
	if err != nil {
		logrus.Println("unable to add to the prodictsmedia table, reason:", err)
		return nil, errors.New("unable to add media")
	}

	urls, err := s.media.Add(data)
	var responseUrls []models.MediaUrl
	for _, v := range keys {
		responseUrls = append(responseUrls, models.MediaUrl{
			Url: urls[v],
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

	urls, err := s.media.Get(ctx, objectIds)
	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	for _, v := range objectIds {
		var media = models.MediaUrl{
			Url: urls[v],
		}
		productMedia = append(productMedia, media)
	}

	return productMedia, nil
}
