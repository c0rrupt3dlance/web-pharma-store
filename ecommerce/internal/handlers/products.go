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

	productResponse, err := h.services.GetById(h.ctx, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	media, err := h.services.FileStorage.GetMedia(h.ctx, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "couldn't get media files"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"product": productResponse,
		"media":   media,
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

	productId, err := h.services.Products.Create(h.ctx, productInput)
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

func (h *Handler) Update(c *gin.Context) {
	var updateProductInput models.UpdateProductInput

	strId := c.Param("id")
	var err error
	productId, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		return
	}

	if err := c.ShouldBindJSON(&updateProductInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid input data",
		})
		return
	}

	err = h.services.Products.Update(h.ctx, productId, updateProductInput)
	if err != nil {
		logrus.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "something went wrong",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{
		"message": "succesfully updated",
	})
}

func (h *Handler) Delete(c *gin.Context) {
	var productId int

	strId := c.Param("id")
	productId, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		return
	}

	err = h.services.Products.Delete(h.ctx, productId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "couldn't delete the product",
		})
	}
}

type categoryIds struct {
	Ids []int `json:"category_ids"`
}

func (h *Handler) GetByCategories(c *gin.Context) {
	var i categoryIds

	if err := c.ShouldBindJSON(&i); err != nil {
		logrus.Printf("user tried to put invalid data%s\n")
		c.JSON(http.StatusBadRequest, gin.H{"message": "invalid data"})
		return
	}

	products, err := h.services.Products.GetByCategories(h.ctx, i.Ids)
	if err != nil {
		logrus.Println("error from service GetByCategories method:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"message": "invalid data"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": products})
}
