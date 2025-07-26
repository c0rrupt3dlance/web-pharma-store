package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
	"strings"
)

const (
	productsTable         = "products"
	categoriesTable       = "categories"
	productsCategoryTable = "products_category"
	productsMediaTable    = "products_media"
)

type ProductPostgres struct {
	pool *pgxpool.Pool
}

func NewProductPostgres(pool *pgxpool.Pool) *ProductPostgres {
	return &ProductPostgres{
		pool: pool,
	}
}

func (r *ProductPostgres) Create(ctx context.Context, p models.ProductInput) (int, error) {
	productQuery := fmt.Sprintf("INSERT INTO %s (name, description, price) values ($1, $2, $3) returning id", productsTable)
	productCategoriesQuery := fmt.Sprintf(`INSERT INTO %s (product_id, category_id) VALUES ($1, $2)`, productsCategoryTable)
	tx, err := r.pool.Begin(ctx)
	logrus.Println(productQuery)
	logrus.Println(productCategoriesQuery)
	if err != nil {

		logrus.Println(err, "Create 1")
		return 0, err
	}
	row := tx.QueryRow(ctx, productQuery, p.Product.Name, p.Product.Description, p.Product.Price)
	if err := row.Scan(&p.Product.Id); err != nil {
		tx.Rollback(ctx)
		logrus.Println(err, "2")
		return 0, err
	}
	for _, i := range p.Categories {
		logrus.Println("categoryId:", i)
		_, err = tx.Exec(ctx, productCategoriesQuery, p.Product.Id, i)
		if err != nil {
			tx.Rollback(ctx)
			logrus.Println(err, "3")
			return 0, err
		}
	}

	tx.Commit(ctx)
	return p.Product.Id, nil
}

func (r *ProductPostgres) GetById(ctx context.Context, productId int) (models.ProductResponse, error) {
	var p = models.ProductResponse{
		Product: models.Product{
			Id: productId,
		},
	}
	query := fmt.Sprintf(`select name, description, price from %s where id=$1`, productsTable)
	categoryQuery := fmt.Sprintf(`SELECT ct.id, ct.name from %s ct inner join
                      %s pc on ct.id = pc.category_id where pc.product_id = $1`, categoriesTable, productsCategoryTable)
	row := r.pool.QueryRow(ctx, query, productId)
	if err := row.Scan(&p.Product.Name, &p.Product.Description, &p.Product.Price); err != nil {
		logrus.Println(err, productId)
		return models.ProductResponse{}, err
	}

	rows, err := r.pool.Query(ctx, categoryQuery, productId)

	if err != nil {
		logrus.Println(err)
		return models.ProductResponse{}, err
	}
	defer rows.Close()
	logrus.Println(rows)
	for rows.Next() {
		category := models.Category{}
		if err = rows.Scan(&category.Id, &category.Name); err != nil {
			logrus.Println(err)
			return models.ProductResponse{}, err
		}
		logrus.Println(p.Categories, category, "categories are")
		p.Categories = append(p.Categories, &category)
	}
	return p, nil
}

func (r *ProductPostgres) Update(ctx context.Context, productId int, p models.UpdateProductInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argsId := 1
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	if p.Name != nil {
		setValues = append(setValues, fmt.Sprintf(`name=$%d`, argsId))
		args = append(args, *p.Name)
		argsId++
	}

	if p.Description != nil {
		setValues = append(setValues, fmt.Sprintf(`description=$%d`, argsId))
		args = append(args, *p.Description)
		argsId++
	}

	if p.Price != nil {
		setValues = append(setValues, fmt.Sprintf(`price=$%d`, argsId))
		args = append(args, *p.Price)
		argsId++
	}
	args = append(args, productId)
	values := strings.Join(setValues, ", ")
	updateProductQuery := fmt.Sprintf(`
			UPDATE %s SET %s WHERE id=$%d
		`, productsTable, values, argsId)
	if len(setValues) > 0 {
		_, err = tx.Exec(ctx, updateProductQuery, args...)
		if err != nil {

			logrus.Printf("error when updating product: %s", err)
			tx.Rollback(ctx)
			return err
		}
	} else {
		tx.Rollback(ctx)
		return errors.New("no fields provided for update")
	}

	if p.Categories != nil {
		newCategoryIds := make(map[int]bool)
		currentCategoriesIds := make(map[int]bool)

		getCategoriesQuery := fmt.Sprintf(`
			SELECT category_id FROM %s WHERE product_id=$1
		`, productsCategoryTable)
		deleteCategoryQuery := fmt.Sprintf(`
			DELETE FROM %s WHERE product_id = $1 and category_id=$2
		`, productsCategoryTable)
		addCategoryQuery := fmt.Sprintf(`
			INSERT INTO %s (product_id, category_id) VALUES ($1, $2)
		`, productsCategoryTable)
		for _, v := range p.Categories {
			newCategoryIds[*v] = true
		}

		rows, err := tx.Query(ctx, getCategoriesQuery, productId)
		if err != nil {
			return err
		}
		defer rows.Close()
		for rows.Next() {
			var currentCategory int
			if err = rows.Scan(&currentCategory); err != nil {
				logrus.Printf("error when getting categories: %s", err)
				return err
			}
			currentCategoriesIds[currentCategory] = true
		}

		for k, _ := range newCategoryIds {
			if _, exists := currentCategoriesIds[k]; !exists {
				_, err = tx.Exec(ctx, addCategoryQuery, productId, k)
				if err != nil {
					logrus.Printf("addCategoryQuery: %s, pr_id: %d, cat_id: %d", addCategoryQuery, productId, k)
					tx.Rollback(ctx)
					return err
				}
			}
		}

		for k, _ := range currentCategoriesIds {
			if _, exists := newCategoryIds[k]; !exists {
				_, err = tx.Exec(ctx, deleteCategoryQuery, productId, k)
				if err != nil {
					tx.Rollback(ctx)
					return err
				}
			}
		}

	}
	if err = tx.Commit(ctx); err != nil {
		return err
	}
	return nil
}
func (r *ProductPostgres) Delete(ctx context.Context, productId int) error {
	query := fmt.Sprintf(`DELETE FROM %s WHERE id=$1`, productsTable)
	_, err := r.pool.Exec(ctx, query, productId)
	if err != nil {
		logrus.Printf("failed to delete product with %d id", productId)
		return err
	}
	return nil
}

func (r *ProductPostgres) GetByCategories(ctx context.Context, categoriesId []int) ([]models.ProductResponse, error) {
	productsByCategory := make([]models.ProductResponse, 0)
	setCategories := make([]string, 0)
	argsId := 1
	for _ = range categoriesId {
		setCategories = append(setCategories, fmt.Sprintf("$%d", argsId))
		argsId++
	}
	categories := strings.Join(setCategories, ", ")
	query := fmt.Sprintf(`
		SELECT pt.* 
		FROM %s pt 
		INNER JOIN %s pc ON pt.id=pc.product_id where pc.category_id IN (%s) 
		GROUP BY pt.id 
		HAVING COUNT(DISTINCT pc.category_id)=%d; 
	`, productsTable, productsCategoryTable, categories, len(categoriesId))

	rows, err := r.pool.Query(ctx, query, intSliceToAnySlice(categoriesId)...)
	if err != nil {
		logrus.Printf(`
			couldn't get rows from the db
			current query: %s
			error from db: %s
			cat_ids: %s
		`, query, err, categoriesId)
		return nil, errors.New("error during getting rows")
	}
	defer rows.Close()
	for rows.Next() {
		var p models.ProductResponse
		if err = rows.Scan(&p.Product.Id, &p.Product.Name, &p.Product.Description, &p.Product.Price); err != nil {
			logrus.Printf(`
			couldn't scan product from the rows
			error from db: %s
			`, err)
			return nil, errors.New("couldn't scan for product")
		}
		productsByCategory = append(productsByCategory, p)
	}
	return productsByCategory, nil
}

func intSliceToAnySlice(s []int) []any {
	anySlice := make([]any, len(s))
	for i, v := range s {
		anySlice[i] = v
	}
	return anySlice
}

func (r *ProductPostgres) AddProductMedia(ctx context.Context, productId int, keys []string) error {
	tx, err := r.pool.Begin(ctx)
	defer tx.Commit(ctx)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`
		INSERT INTO %s(product_id, media_id) VALUES($1, $2)
	`, productsMediaTable)

	for _, v := range keys {
		_, err = tx.Exec(ctx, query, productId, v)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	logrus.Println(productId, "got new media files")
	return nil
}
