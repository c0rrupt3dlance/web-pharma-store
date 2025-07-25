package main

import (
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/app"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/handlers"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/repository"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/services"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"os"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		logrus.Printf("failed to load .env due to: %s", err)
		os.Exit(1)
	}
	pool, err := repository.NewPgPool(
		repository.PgConfig{
			Host:     os.Getenv("POSTGRES_HOST"),
			Port:     os.Getenv("POSTGRES_PORT"),
			User:     os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
			Database: os.Getenv("POSTGRES_DB"),
		},
		context.Background(),
	)
	if err != nil {
		log.Println("failed to create postgres pool:", err)
		os.Exit(1)
	}

	minioClient, err := repository.NewMinioClient(
		context.Background(),
		repository.MediaConfig{
			os.Getenv("MINIO_ENDPOINT"),
			os.Getenv("MINIO_ACCESS_KEY"),
			os.Getenv("MINIO_SECRET_KEY"),
			os.Getenv("MINIO_BUCKET"),
			true,
		},
	)

	var repo = repository.NewRepository(pool)
	var service = services.NewService(repo, os.Getenv("SIGNING_KEY"))
	var handler = handlers.NewHandler(context.Background(), service)
	server := new(app.Server)

	err = server.Run(os.Getenv("SERVER_PORT"), handler.InitRoutes())
	if err != nil {
		logrus.Println("failed to start server due to: %s", err)
		os.Exit(1)
	}
}
