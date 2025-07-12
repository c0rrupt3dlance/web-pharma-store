package handlers

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strconv"
)

func (h *Handler) GetById(c *gin.Context) {
	productId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		return
	}

	productResponse, err := h.services.GetById(productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": productResponse,
	})
}
func (h *Handler) Create(c *gin.Context) {
	var productInput models.ProductInput
	if err := c.ShouldBindJSON(&productInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input data",
		})
		return
	}

	logrus.Println(productInput)

	productId, err := h.services.Create(productInput)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": productId,
	})
}
func (h *Handler) Delete(c *gin.Context) {}
func (h *Handler) Update(c *gin.Context) {}
