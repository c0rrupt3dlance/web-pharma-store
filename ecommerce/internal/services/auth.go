package services

import (
	"errors"
	"github.com/c0rrupt3dlance/web-pharma-store/ecommerce/internal/models"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	signingKey string
}

type tokenClaims struct {
	UserId   int
	Username string
	Role     string
	jwt.RegisteredClaims
}

func NewAuthService(signingKey string) *AuthService {
	return &AuthService{
		signingKey: signingKey,
	}
}

func (s *AuthService) VerifyAccessToken(tokenString string) (models.User, error) {
	user := models.User{}
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			logrus.Println("invalid sign")
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.signingKey), nil
	})

	if err != nil {
		return models.User{}, err
	}

	claims, ok := token.Claims.(*tokenClaims)
	if ok && token.Valid {
		user.Id = claims.UserId
		user.Username = claims.Username
		user.Role = claims.Role
		return user, nil
	}

	return models.User{}, errors.New("tokens claims are invalid or token is just invalid")
}
