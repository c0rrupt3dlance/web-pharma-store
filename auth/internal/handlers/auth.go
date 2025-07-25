package handlers

import (
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (h *Handler) SignUp(c *gin.Context) {
	var user models.User

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

type SignInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *Handler) SignIn(c *gin.Context) {
	var input SignInInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "bad input"})
		return
	}

	accessToken, refreshToken, err := h.services.GenerateTokens(input.Username, input.Password)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid credentials"})
		return
	}
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    refreshToken,
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	},
	)
	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}

type RefreshInput struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

func (h *Handler) RefreshToken(c *gin.Context) {
	var input RefreshInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token"})
		return
	}

	accessToken, newRefreshToken, err := h.services.RefreshTokens(input.RefreshToken)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid refresh token"})
		return
	}

	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "refresh_token",
		Value:    newRefreshToken,
		MaxAge:   7 * 24 * 60 * 60,
		Path:     "/",
		Domain:   "",
		Secure:   false,
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
	})

	c.JSON(http.StatusOK, gin.H{
		"access_token": accessToken,
	})
}
