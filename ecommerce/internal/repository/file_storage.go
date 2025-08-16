package repository

import (
	"bytes"
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/minio/minio-go/v7"
	"sync"
	"time"
)

const numWorkers = 5

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

type FileChannel struct {
	ObjectId string
	File     models.FileDataType
}

func (r *MinioFileStorage) worker(jobs <-chan FileChannel, urlCh chan<- models.MediaUrl, wg *sync.WaitGroup, ctx context.Context, cancel context.CancelFunc) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}
			_, err := r.Client.PutObject(ctx, r.Bucket, job.ObjectId, bytes.NewReader(job.File.Data),
				int64(len(job.File.Data)), minio.PutObjectOptions{
					ContentType: job.File.DataType,
				})
			if err != nil {
				cancel()
				return
			}

			Url, err := r.Client.PresignedGetObject(ctx, r.Bucket, job.ObjectId, time.Hour*1, nil)
			if err != nil {
				cancel()
				return
			}

			urlCh <- models.MediaUrl{
				ObjectId: job.ObjectId,
				Url:      Url.String(),
			}
		}
	}
}

func (r *MinioFileStorage) AddMedia(data map[string]models.FileDataType) (map[string]string, error) {
	urls := make(map[string]string)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wg := sync.WaitGroup{}

	jobs := make(chan FileChannel, len(data))
	urlCh := make(chan models.MediaUrl, len(data))

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go r.worker(jobs, urlCh, &wg, ctx, cancel)
	}

	for objectId, file := range data {
		jobs <- FileChannel{ObjectId: objectId, File: file}
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(urlCh)
	}()

	for url := range urlCh {
		urls[url.ObjectId] = url.Url
	}

	return urls, nil
}

func (r *MinioFileStorage) GetMedia(ctx context.Context, objectIds []string) (map[string]string, error) {
	urls := make(map[string]string)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	urlCh := make(chan models.MediaUrl, len(objectIds))
	var wg sync.WaitGroup

	for _, v := range objectIds {
		wg.Add(1)

		go func(objectId string) {
			defer wg.Done()

			Url, err := r.Client.PresignedGetObject(ctx, r.Bucket, objectId, time.Hour*1, nil)
			if err != nil {
				cancel()
				return
			}
			urlCh <- models.MediaUrl{
				ObjectId: v,
				Url:      Url.String(),
			}
		}(v)

	}

	go func() {
		wg.Wait()
		close(urlCh)
	}()

	for link := range urlCh {
		urls[link.ObjectId] = link.Url
	}

	return urls, nil

}
