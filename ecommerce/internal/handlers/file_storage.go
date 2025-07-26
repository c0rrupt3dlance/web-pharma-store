package handlers

import (
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strconv"
)

func (h *Handler) Upload(c *gin.Context) {
	strId := c.Param("id")
	productId, err := strconv.Atoi(strId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid product id",
		})
		return
	}

	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid multiform",
		})
		return
	}
	files := form.File["media"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "no file was uploaded",
		})
		return
	}
	var mediaFiles []models.FileDataType

	for _, fileHeader := range files {
		file, err := fileHeader.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "cannot open uploaded file",
			})
			return
		}
		defer file.Close()

		data, err := io.ReadAll(file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "cannot open uploaded file",
			})
			return
		}
		mediaFiles = append(mediaFiles, models.FileDataType{
			FileName: fileHeader.Filename,
			Data:     data,
			DataType: fileHeader.Header.Get("Content-Type"),
		})
	}

	urls, err := h.services.FileStorage.AddMedia(h.ctx, productId, mediaFiles)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "got some problems",
		})
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"data": urls})
}
