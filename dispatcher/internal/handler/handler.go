package handler

import (
	"context"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *service.Service
	ctx     context.Context
}

func NewHandler(ctx context.Context, service *service.Service) *Handler {
	return &Handler{
		service,
		ctx,
	}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.Use(
		gin.Logger(),
		gin.Recovery(),
	)

	api := router.Group("/api/v1")
	{
		orders := api.Group("/dispatcher")
		{
			orders.GET("/get-orders")    // expect for it to return all the existing dispatcher or by some parameter
			orders.POST("/create-order") // expected to be called somehow only from api gateway
		}
	}
}
