package repository

import (
	"bytes"
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/minio/minio-go/v7"
	"sync"
	"time"
)

type MinioFileStorage struct {
	Client *minio.Client
	Bucket string
}

func NewMinioFileStorage(client *minio.Client, bucket string) *MinioFileStorage {
	return &MinioFileStorage{
		Client: client,
		Bucket: bucket,
	}
}

func (r *MinioFileStorage) Add(data map[string]models.FileDataType) (map[string]string, error) {
	urls := make(map[string]string, len(data))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	urlCh := make(chan models.MediaUrl, len(data))

	var wg sync.WaitGroup

	for objectId, file := range data {
		wg.Add(1)
		go func(objectId string, file models.FileDataType) {
			defer wg.Done()
			_, err := r.Client.PutObject(ctx, r.Bucket, objectId, bytes.NewReader(file.Data),
				int64(len(file.Data)), minio.PutObjectOptions{})
			if err != nil {
				cancel()
				return
			}

			Url, err := r.Client.PresignedGetObject(ctx, r.Bucket, objectId, time.Second*24*60*60, nil)
			if err != nil {
				cancel()
				return
			}

			urlCh <- models.MediaUrl{
				ObjectId: objectId,
				Url:      Url.String(),
			}
		}(objectId, file)

	}

	go func() {
		wg.Wait()
		close(urlCh)
	}()

	for link := range urlCh {
		urls[link.ObjectId] = link.Url
	}

	_ = len(urls)
	return urls, nil
}
