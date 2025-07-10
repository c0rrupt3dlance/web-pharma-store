package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) GetById(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "handlers work",
	})
}
func (h *Handler) Create(c *gin.Context) {}
func (h *Handler) Delete(c *gin.Context) {}
func (h *Handler) Update(c *gin.Context) {}
