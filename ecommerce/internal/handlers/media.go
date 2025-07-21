package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) AddMedia(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "bad request",
		})
		return
	}

	files := form.File["files"]
	if files == nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "no images or video provided"})
		return
	}
}
