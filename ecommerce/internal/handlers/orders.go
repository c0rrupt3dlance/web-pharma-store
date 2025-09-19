package handlers

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

func (h *Handler) CreateOrder(c *gin.Context) {
	var (
		input models.OrderInput
		err   error
	)
	userId, ok := c.Get(userCtx)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "invalid user token",
		})
		return
	}

	input.ClientId = userId.(int)
	if err = c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid data",
		})
		return
	}

	orderId, err := h.services.Orders.CreateOrder(h.ctx, input)
	if err != nil {
		logrus.Printf("error from service: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message":  "order created",
		"order_id": orderId,
	})
}

func (h *Handler) GetOrder(c *gin.Context) {
}

func (h *Handler) GetAllOrders(c *gin.Context) {
}
