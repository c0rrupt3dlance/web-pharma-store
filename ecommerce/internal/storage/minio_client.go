package storage

import (
	"context"

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
