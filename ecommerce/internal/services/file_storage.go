package services

import (
	"context"
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
	objects := make(map[int]string, len(mediaFiles))
	files := make(map[string]models.FileDataType, len(mediaFiles))
	for _, v := range mediaFiles {
		id := uuid.New().String()
		objects[v.Position] = id
		files[id] = v
	}

	err := s.products.AddProductMedia(ctx, productId, objects)
	if err != nil {
		return nil, err
	}

	urls, err := s.media.AddMedia(ctx, files)
	if err != nil {
		return nil, err
	}

	var responseUrls []models.MediaUrl
	for k, v := range objects {
		responseUrls = append(responseUrls, models.MediaUrl{
			Url:      urls[v],
			Position: k,
		})
	}

	return responseUrls, nil
}

func (s *FileStorageService) GetMedia(ctx context.Context, productId int) ([]models.MediaUrl, error) {
	objectIds, err := s.products.GetProductMedia(ctx, productId)
	if err != nil {
		logrus.Println(err)
		return nil, err
	}

	urls, err := s.media.GetMedia(ctx, objectIds)
	if err != nil {
		logrus.Println(err)
		return nil, fmt.Errorf("error while %w: ", err)
	}

	for i, _ := range objectIds {
		objectIds[i].Url = urls[objectIds[i].ObjectId]
	}

	return objectIds, nil
}
