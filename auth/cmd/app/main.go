package main

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/app"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/handlers"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/services"
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
	)
	if err != nil {
		log.Println("failed to create postgres pool")
		os.Exit(1)
	}
	repo := repository.NewRepository(pool)
	deps := services.Dependencies{
		Repo:       repo,
		SigningKey: os.Getenv("SIGNING_KEY"),
	}
	service := services.NewService(deps.Repo, deps.SigningKey)
	handler := handlers.NewHandler(service)

	server := new(app.Server)

	err = server.Run(os.Getenv("SERVER_PORT"), handler.InitRoutes())
	if err != nil {
		logrus.Println("failed to start server due to: %s", err)
		os.Exit(1)
	}
}
