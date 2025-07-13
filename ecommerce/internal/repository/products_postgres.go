package repository

import (
	"context"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

const (
	productsTable         = "products"
	categoriesTable       = "categories"
	productsCategoryTable = "products_category"
)

type ProductPostgres struct {
	pool *pgxpool.Pool
}

func NewProductPostgres(pool *pgxpool.Pool) *ProductPostgres {
	return &ProductPostgres{
		pool: pool,
	}
}

func (r *ProductPostgres) Create(p models.ProductInput) (int, error) {
	productQuery := fmt.Sprintf("INSERT INTO %s (name, description, price) values ($1, $2, $3) returning id", productsTable)
	productCategoriesQuery := fmt.Sprintf(`INSERT INTO %s (product_id, category_id) VALUES ($1, $2)`, productsCategoryTable)
	tx, err := r.pool.Begin(context.Background())
	logrus.Println(productQuery)
	logrus.Println(productCategoriesQuery)
	if err != nil {

		logrus.Println(err, "Create 1")
		return 0, err
	}
	row := tx.QueryRow(context.Background(), productQuery, p.Product.Name, p.Product.Description, p.Product.Price)
	if err := row.Scan(&p.Product.Id); err != nil {
		tx.Rollback(context.Background())
		logrus.Println(err, "2")
		return 0, err
	}
	for _, i := range p.Categories {
		logrus.Println("categoryId:", i)
		_, err = tx.Exec(context.Background(), productCategoriesQuery, p.Product.Id, i)
		if err != nil {
			tx.Rollback(context.Background())
			logrus.Println(err, "3")
			return 0, err
		}
	}

	tx.Commit(context.Background())
	return p.Product.Id, nil
}

func (r *ProductPostgres) GetById(ProductId int) (models.ProductResponse, error) {
	var p models.ProductResponse
	query := fmt.Sprintf(`select * from %s where id=$1`, productsTable)
	categoryQuery := fmt.Sprintf(`SELECT ct.id, ct.name from %s ct inner join
                      %s pc on ct.id = pc.category_id where pc.product_id = $1`, categoriesTable, productsCategoryTable)
	row := r.pool.QueryRow(context.Background(), query, ProductId)
	if err := row.Scan(&p.Product.Id, &p.Product.Name, &p.Product.Description, &p.Product.Price); err != nil {
		logrus.Println(err, "1", ProductId)
		return models.ProductResponse{}, err
	}

	rows, err := r.pool.Query(context.Background(), categoryQuery, ProductId)

	if err != nil {
		logrus.Println(err, "2")
		return models.ProductResponse{}, err
	}
	defer rows.Close()
	logrus.Println(rows)
	for rows.Next() {
		category := models.Category{}
		if err = rows.Scan(&category.Id, &category.Name); err != nil {
			logrus.Println(err, "3")
			return models.ProductResponse{}, err
		}
		logrus.Println(p.Categories, category, "categories are")
		p.Categories = append(p.Categories, category)
	}
	return p, nil
}
func (r *ProductPostgres) Update(p models.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1

	if p.Name != nil {
		setValues = append(setValues, fmt.Sprintf(`name=$%d`, argsId))
		args = append(args, *p.Name)
		argsId++
	}

	if p.Description != nil {
		setValues = append(setValues, fmt.Sprintf(`name=$%d`, argsId))
		args = append(args, *p.Description)
		argsId++
	}

	if p.Price != nil {
		setValues = append(setValues, fmt.Sprintf(`name=$%d`, argsId))
		args = append(args, *p.Price)
		argsId++
	}

	if p.Categories != nil {
		newCategoryIds := make(map[int]bool)
		getCategoriesQuery := fmt.Sprintf(`
			get id from %s ct inner join %s pc 
			on ct.id=pc.category_id where pc.product_id=$1
		`)
		deleteCateforyQuery := fmt.Sprintf(`
			
		`)
		for _, v := range p.Categories {
			newCategoryIds[*v] = true
		}
		rows, err := r.pool.Query(context.Background(), getCategoriesQuery, p.Id)
		if err != nil {
			return err
		}
		defer rows.Close()
		currentCategoriesIds := make(map[int]bool)
		for rows.Next() {
			var currentCategory int
			if err = rows.Scan(&currentCategory); err != nil {
				return err
			}
			currentCategoriesIds[currentCategory] = true
		}

		for k, _ := range newCategoryIds {
			if _, exists := currentCategoriesIds[k]; !exists {

			}
		}
	}
	return nil
}
func (r *ProductPostgres) Delete(ProductId int) error {
	return nil
}
