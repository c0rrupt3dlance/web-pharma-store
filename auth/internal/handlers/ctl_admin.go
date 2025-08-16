package handlers

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) RegisterAdmin(c *gin.Context) {
	var user = models.User{
		Role: "admin",
	}

	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId, err := h.services.Create(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"user_id": userId})
}
