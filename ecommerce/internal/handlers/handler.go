package handlers

import (
	"context"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Service
	ctx      context.Context
}

func NewHandler(ctx context.Context, services *services.Service) *Handler {
	return &Handler{
		services: services,
		ctx:      ctx,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(
		gin.Logger(), gin.Recovery(),
	)
	router.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ping": "pong",
		})
	})
	api := router.Group("/api/v1", h.verifyUserToken)
	{
		products := api.Group("/products")
		{
			products.GET("/:id", h.GetById)
			products.POST("/get_by_category", h.GetByCategories)
			products.GET("/cart", h.GetUserCart)
			products.POST("/:id/add_to_cart", h.AddProductToCart)
			products.DELETE("/cart/:id", h.DeleteProductFromCart)

			ctlStore := products.Group("/ctl", h.verifyAdminToken)
			{
				ctlStore.POST("/", h.Create)
				ctlStore.PUT("/:id", h.Update)
				ctlStore.DELETE("/:id", h.Delete)
				ctlStore.POST("/:id/media", h.Upload)
			}

		}

		orders := api.Group("/orders")
		{
			orders.POST("/", h.CreateOrder)
			orders.GET("/", h.GetAllOrders)
			orders.GET("/:id", h.GetOrder)
		}

		return router
	}
}
