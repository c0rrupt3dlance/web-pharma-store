package repository

import (
	"bytes"
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"sync"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/sirupsen/logrus"
)

type MinioClient struct {
	Client *minio.Client
	Bucket string
}

func NewMinioClient(ctx context.Context, endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		logrus.Println("error when creating minio client")
		return nil, err
	}

	exists, err := client.BucketExists(ctx, bucket)
	if err != nil {
		logrus.Println("error during checking if bucket exists")
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, bucket, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Println("error when creating bucket")
			return nil, err
		}
		logrus.Printf("bucket %s created successfully", bucket)
	} else {
		logrus.Printf("bucket %s already exists", bucket)
	}

	return &MinioClient{
		Client: client,
		Bucket: bucket,
	}, nil
}

func (r *MinioClient) Create(data map[string]models.FileDataType) (map[string]string, error) {
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
