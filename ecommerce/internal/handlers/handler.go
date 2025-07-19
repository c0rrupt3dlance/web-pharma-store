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
	api := router.Group("/api/v1", h.parseAccessToken)
	{
		products := api.Group("/products")
		{
			products.GET("/:id", h.GetById)
			products.POST("/get_by_category", h.GetByCategories)
			products.POST("/", h.Create)
			products.PUT("/:id", h.Update)
			products.DELETE("/:id", h.Delete)
			products.GET("/cart")
			products.POST("/:id/add_to_cart")
			products.DELETE("/cart/")
		}

	}

	return router
}
