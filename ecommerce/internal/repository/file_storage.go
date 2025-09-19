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
	Client       *minio.Client
	Bucket       string
	WorkerShards int
}

type fileObject struct {
	objectId string
	file     models.FileDataType
}

func NewMinioFileStorage(client *minio.Client, bucket string, shards int) *MinioFileStorage {
	return &MinioFileStorage{
		Client:       client,
		Bucket:       bucket,
		WorkerShards: shards,
	}
}

func (r *MinioFileStorage) AddMedia(parentCtx context.Context, data map[string]models.FileDataType) (map[string]string, error) {
	objects := make(map[string]string, len(data))
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	files := make(chan fileObject, len(data))
	urlCh := make(chan models.MediaUrl, len(data))

	var wg sync.WaitGroup

	for _ = range r.WorkerShards {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-files:
					if !ok {
						return
					}
					_, err := r.Client.PutObject(ctx, r.Bucket, job.objectId, bytes.NewReader(job.file.Data),
						int64(len(job.file.Data)), minio.PutObjectOptions{
							ContentType: job.file.DataType,
						})

					if err != nil {
						return
					}

					url, err := r.Client.PresignedGetObject(ctx, r.Bucket, job.objectId, time.Hour, nil)
					if err != nil {
						return
					}

					urlCh <- models.MediaUrl{
						ObjectId: job.objectId,
						Url:      url.String(),
					}
				}
			}
		}()
	}

	for objId, file := range data {
		files <- fileObject{
			objectId: objId,
			file:     file,
		}
	}
	close(files)

	go func() {
		wg.Wait()
		close(urlCh)
	}()

	for url := range urlCh {
		objects[url.ObjectId] = url.Url
	}

	return objects, nil
}

func (r *MinioFileStorage) GetMedia(parentCtx context.Context, objectIds []models.MediaUrl) (map[string]string, error) {
	objects := make(map[string]string, len(objectIds))
	jobs := make(chan string, len(objectIds))
	urls := make(chan models.MediaUrl, len(objectIds))
	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	wg := sync.WaitGroup{}
	for _ = range r.WorkerShards {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				select {
				case <-ctx.Done():
					return
				case job, ok := <-jobs:
					if !ok {
						return
					}
					url, err := r.Client.PresignedGetObject(ctx, r.Bucket, job, time.Hour, nil)
					if err != nil {
						return
					}

					urls <- models.MediaUrl{
						ObjectId: job,
						Url:      url.String(),
					}
				}
			}
		}()
	}

	for _, id := range objectIds {
		jobs <- id.ObjectId
	}
	close(jobs)

	go func() {
		wg.Wait()
		close(urls)
	}()

	for url := range urls {
		objects[url.ObjectId] = url.Url
	}

	return objects, nil
}
