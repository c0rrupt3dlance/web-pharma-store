package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type OrderPostgres struct {
	pool *pgxpool.Pool
}

var Statuses = map[string]int{
	"created":         1,
	"being delivered": 2,
	"delivered":       3,
	"cancelled":       4,
	"returned":        5,
}

func NewOrderPostgresRepo(pool *pgxpool.Pool) *OrderPostgres {
	return &OrderPostgres{
		pool,
	}
}

func (r *OrderPostgres) CreateOrder(ctx context.Context, input models.OrderInput) (int, error) {
	tx, err := r.pool.Begin(ctx)

	if err != nil {
		return 0, fmt.Errorf("couldn't create a transaction %s", err)
	}

	defer func() {
		if err != nil {
			tx.Rollback(ctx)
		}
	}()

	var receiverId int

	createReceiverQuery := fmt.Sprintf(`
	INSERT INTO %s(name,phone,address)
	VALUES ($1, $2, $3) 
	RETURNING id
	`, receiversTable)

	createOrderQuery := fmt.Sprintf(`
	INSERT INTO %s(client_id, receiver_id, status_id)
	VALUES ($1, $2, $3) 
	RETURNING id
	`, ordersTable)

	selectProductQuery := fmt.Sprintf(`
	SELECT name, price FROM %s WHERE id=$1
	`, productsTable)

	createItemQuery := fmt.Sprintf(`
	INSERT INTO %s(order_id, item_id, item_name, quantity, price)
	VALUES($1, $2, $3, $4, $5)
	`, orderItemsTable)

	row := tx.QueryRow(ctx, createReceiverQuery,
		input.Receiver.Name, input.Receiver.Phone, input.Receiver.Address)

	if err = row.Scan(&receiverId); err != nil {
		return 0, fmt.Errorf("error: %s", err)
	}

	row = tx.QueryRow(ctx, createOrderQuery,
		input.ClientId, receiverId, Statuses["created"])

	var orderId int
	if err = row.Scan(&orderId); err != nil {
		return 0, fmt.Errorf("error: %s", err)
	}

	for _, v := range input.Items {
		var (
			itemName  string
			itemPrice int
		)

		row = tx.QueryRow(ctx, selectProductQuery, v.ItemId)

		if err = row.Scan(&itemName, &itemPrice); err != nil {
			return 0, fmt.Errorf("error: %s", err)
		}

		_, err = tx.Exec(ctx, createItemQuery,
			orderId, v.ItemId, itemName, v.Quantity, itemPrice)

		if err != nil {
			return 0, fmt.Errorf("error: %s", err)
		}
	}

	buildPayloadQuery := fmt.Sprintf(`
	SELECT json_build_object(
	'id', o.id,
	'client_id', o.client_id,
	'receiver', json_build_object(
		'id', r.id,
		'name', r.name,
		'phone', r.phone,
		'address', r.address
	),
	'status', json_build_object(
		'id', os.id,
		'code', os.code
	),
	'items', COALESCE(
	         	json_agg(
	         		json_build_object(
	         			'id', oi.id,
	         			'item_id', oi.item_id,
	         			'item_name', oi.item_name,
	         			'quantity', oi.quantity,
	         			'price', oi.price
	         		)
	         ) FILTER (WHERE oi.id IS NOT NULL),
	         '[]'
	         ),
	'created_at', o.created_at,
	'update_at', o.updated_at
	) as payload FROM %s o
	JOIN %s r on r.id = o.receiver_id
	JOIN %s os on os.id = o.status_id
	LEFT JOIN %s oi on oi.order_id = o.id
	WHERE o.id=$1
	GROUP BY o.id, r.id, os.id
	`, ordersTable, receiversTable, orderStatusesTable, orderItemsTable)

	var payload []byte
	if err = tx.QueryRow(ctx, buildPayloadQuery, orderId).Scan(&payload); err != nil {
		logrus.Errorf("error %s", err)
		return 0, errors.New("couldn't get order payload")
	}

	AddToOutBoxQuery := fmt.Sprintf(`
	INSERT INTO %s (aggregate_type, aggregate_id, event_type, 
	                   payload)
	VALUES ($1, $2, $3, $4)
	`, outboxTable)

	_, err = tx.Exec(ctx, AddToOutBoxQuery, "orders", orderId, "order_created", payload)
	if err != nil {
		return 0, err
	}

	if err = tx.Commit(ctx); err != nil {
		return 0, err
	}
	return orderId, nil
}

func (r *OrderPostgres) GetOrder(ctx context.Context, ClientId int) (models.Order, error) {
	return models.Order{}, nil
}

func (r *OrderPostgres) GetAllOrders(ctx context.Context, ClientId int) ([]models.Order, error) {
	return make([]models.Order, 0), nil
}

func (r *OrderPostgres) AddToOutbox(ctx context.Context, orderId int) error {
	return nil
}
