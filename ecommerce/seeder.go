package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"log"
	"math/rand"
	"os"
	"time"
)

func main() {
	err := godotenv.Load(".env")
	ctx := context.Background()
	pool, err := pgxpool.New(ctx, fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_HOST"),
		os.Getenv("POSTGRES_PORT"),
		os.Getenv("POSTGRES_DB")))

	if err != nil {
		logrus.Printf("fail due to %s", err)
		os.Exit(1)
	}

	err = pool.Ping(ctx)
	if err != nil {
		os.Exit(1)
	}

	pool.Exec(ctx, `DELETE FROM products_category`)
	pool.Exec(ctx, `DELETE FROM products`)
	pool.Exec(ctx, `DELETE FROM categories`)

	// Сидим категории
	categories := []string{"Painkillers", "Antibiotics", "Vitamins", "Antivirus"}
	categoryIDs := []int{}

	for _, name := range categories {
		var id int
		err := pool.QueryRow(ctx, `INSERT INTO categories (name) VALUES ($1) RETURNING id`, name).Scan(&id)
		if err != nil {
			log.Fatalf("Error when seeding categories: %v", err)
		}
		categoryIDs = append(categoryIDs, id)
	}
	products := []struct {
		name        string
		description string
		price       float64
	}{
		{"Paracetamol", "To temp", 120.50},
		{"Fgere", "For pain", 340.00},
		{"Someticin", "Someticin", 250.99},
		{"Ibuprofen", "Saves in tarkov", 180.00},
		{"Arbdol", "Antivitus", 560.00},
	}

	rand.Seed(time.Now().UnixNano())

	for _, p := range products {
		var prodID int
		err := pool.QueryRow(
			ctx, `INSERT INTO products (name, description, price) VALUES ($1, $2, $3) RETURNING id`,
			p.name, p.description, p.price,
		).Scan(&prodID)
		if err != nil {
			log.Fatalf("Error when seeding products: %v", err)
		}
		rand.Shuffle(len(categoryIDs), func(i, j int) {
			categoryIDs[i], categoryIDs[j] = categoryIDs[j], categoryIDs[i]
		})

		for i := 0; i < rand.Intn(2)+1; i++ {
			_, err := pool.Exec(ctx, `INSERT INTO products_category (product_id, category_id) VALUES ($1, $2)`, prodID, categoryIDs[i])
			if err != nil {
				log.Printf("Ошибка вставки в products_category: %v", err)
			}
		}
	}

	fmt.Println("Finished Seeding")
}
