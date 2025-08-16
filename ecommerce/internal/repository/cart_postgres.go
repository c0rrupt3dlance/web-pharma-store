package repository

import (
	"context"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type CartPostgres struct {
	pool *pgxpool.Pool
}

func NewCartPostgres(pool *pgxpool.Pool) *CartPostgres {
	return &CartPostgres{
		pool,
	}
}

func (r *CartPostgres) GetUserCart(ctx context.Context, userId int) (*models.UserCart, error) {
	var userCart = models.UserCart{
		UserId:    userId,
		CartItems: make([]models.CartItem, 0),
	}
	query := fmt.Sprintf(`
	SELECT ci.id,ci.product_id,ci.quantity,pt.price
	FROM %s ci 
	INNER JOIN %s pt
	ON ci.product_id=pt.id
	WHERE ci.user_id=$1
	`, cartItemsTable, productsTable)
	rows, err := r.pool.Query(ctx, query, userId)
	if err != nil {
		logrus.Println("GetUserCart:", err)
		return nil, err
	}
	for rows.Next() {
		var cartItem = models.CartItem{}
		if err := rows.Scan(&cartItem.Id, &cartItem.ProductId, &cartItem.Quantity, &cartItem.Price); err != nil {
			logrus.Println("GetUserCart:", err)
			return nil, err
		}
		userCart.CartItems = append(userCart.CartItems, cartItem)
	}
	return nil, nil
}

func (r *CartPostgres) AddToCart(ctx context.Context, userId int, itemInput models.CartItem) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		logrus.Println("AddToCart", err)
		return err
	}
	query := fmt.Sprintf(`
		INSERT INTO %s (user_id, product_id, quantity, price) 
		VALUES($1,$2,$3,$4)
	`, cartItemsTable)

	_, err = tx.Exec(ctx, query, userId, itemInput)
	if err != nil {
		logrus.Println("AddToCart", err)
		tx.Rollback(ctx)
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		logrus.Println("AddToCart", err)
		return err
	}
	return nil
}

func (r *CartPostgres) RemoveFromCart(ctx context.Context, userId, cartItemId int) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		logrus.Println("RemoveFromCart: ", err)
		return err
	}

	query := fmt.Sprintf(`
	DELETE FROM %s WHERE id=$1 AND user_id=$2 
	`, cartItemsTable)

	_, err = tx.Exec(ctx, query, cartItemId, userId)
	if err != nil {
		logrus.Println("RemoveFromCart: ", err)
		return err
	}

	tx.Commit(ctx)
	return nil
}
