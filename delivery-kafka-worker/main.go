package main

import (
	"context"
	"delivery-kafka-worker/producer"
	"delivery-kafka-worker/repository"
	"delivery-kafka-worker/worker"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"time"
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
		logrus.Println(err)
		os.Exit(1)
	}
	queueRepo := repository.NewMQPostgres(pool)
	queueProducer := producer.NewProducer([]string{"localhost:9092"})
	queueWorker := worker.NewWorker(queueRepo, queueProducer, 2*time.Second)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	queueWorker.Start(ctx)
}
