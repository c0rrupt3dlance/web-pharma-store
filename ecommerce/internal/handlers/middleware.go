package handlers

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) verifyUserToken(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "no access token",
		})
		return
	}
	var tokenParts = strings.Split(header, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		logrus.Println(tokenParts)
		logrus.Println("Twin we got wrong token twin")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	userId, userRole, err := h.services.VerifyAccessToken(tokenParts[1])
	if err != nil || userRole == nil {
		logrus.Printf("%s is our error on service layer", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	c.Set(userCtx, userId)
	c.Next()
}

func (h *Handler) verifyAdminToken(c *gin.Context) {
	header := c.GetHeader(authorizationHeader)
	if header == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "no access token",
		})
		return
	}
	var tokenParts = strings.Split(header, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		logrus.Println(tokenParts)
		logrus.Println("invalid token")
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	userId, userRole, err := h.services.VerifyAccessToken(tokenParts[1])
	if err != nil {
		logrus.Printf("%s is our error on service layer", err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}
	if *userRole != "admin" {
		logrus.Printf("%s have tried to access ctl-panel", userId)
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"message": "invalid token",
		})
		return
	}

	c.Set(userCtx, userId)
	c.Next()
}
