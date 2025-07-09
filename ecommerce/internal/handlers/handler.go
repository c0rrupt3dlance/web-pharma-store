package handlers

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/services"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *services.Service
}

func NewHandler(services *services.Service) *Handler {
	return &Handler{
		services: services,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"ping": "pong",
		})
	})
	api := router.Group("/api/v1", h.ParseAccessToken)
	{
		products := api.Group("/products")
		{
			products.GET("/:id")
			products.POST("/")
			products.PUT("/:id")
			products.DELETE("/:id")
		}
	}

	return router
}
