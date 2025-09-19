package repository

import (
	"context"
	models "github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/minio/minio-go/v7"
)

type Products interface {
	Create(ctx context.Context, product models.ProductInput) (int, error)
	GetById(ctx context.Context, productId int) (models.ProductResponse, error)
	Update(ctx context.Context, productId int, product models.UpdateProductInput) error
	Delete(ctx context.Context, ProductId int) error
	GetByCategories(ctx context.Context, categoriesId []int) ([]models.ProductResponse, error)
	AddProductMedia(ctx context.Context, productId int, objIds map[int]string) error
	GetProductMedia(ctx context.Context, productId int) ([]models.MediaUrl, error)
}

type Cart interface {
	AddItem(userId, productId, quantity int) error
	UpdateQuantity(userId, productId, quantity int) error
	RemoveItem(userId, productId int) error
	GetCart(userId int) ([]models.CartItem, error)
	ClearCart(userId int) error
}

type Orders interface {
	CreateOrder(ctx context.Context, input models.OrderInput) (int, error)
	GetOrder(ctx context.Context, clientId int) (models.Order, error)
	GetAllOrders(ctx context.Context, clientId int) ([]models.Order, error)
	AddToOutbox(ctx context.Context, orderId int) error
}

type FileStorage interface {
	AddMedia(ctx context.Context, data map[string]models.FileDataType) (map[string]string, error)
	GetMedia(ctx context.Context, objectIds []models.MediaUrl) (map[string]string, error)
}
type Repository struct {
	Products
	Cart
	FileStorage
	Orders
}

func NewRepository(pool *pgxpool.Pool, client *minio.Client, bucket string, shards int) *Repository {
	return &Repository{
		Products:    NewProductPostgres(pool),
		FileStorage: NewMinioFileStorage(client, bucket, shards),
		Orders:      NewOrderPostgresRepo(pool),
	}
}
