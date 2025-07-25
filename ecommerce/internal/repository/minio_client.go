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

type MediaConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func NewMinioClient(ctx context.Context, cfg MediaConfig) (*minio.Client, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		logrus.Println("error when creating minio client")
		return nil, err
	}

	exists, err := client.BucketExists(ctx, cfg.Bucket)
	if err != nil {
		logrus.Println("error during checking if Bucket exists")
		return nil, err
	}

	if !exists {
		err = client.MakeBucket(ctx, cfg.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			logrus.Println("error when creating Bucket")
			return nil, err
		}
		logrus.Printf("Bucket %s created successfully", cfg.Bucket)
	} else {
		logrus.Printf("Bucket %s already exists", cfg.Bucket)
	}

	return client, nil
}
