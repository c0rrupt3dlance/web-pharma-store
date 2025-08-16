package services

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type AuthService struct {
	signingKey string
}

type tokenClaims struct {
	userId   int
	username string
	role     string
	jwt.RegisteredClaims
}

func NewAuthService(signingKey string) *AuthService {
	return &AuthService{
		signingKey: signingKey,
	}
}

func (s *AuthService) VerifyAccessToken(tokenString string) (int, *string, error) {
	var userId int
	var userRole string

	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}
		return []byte(s.signingKey), nil
	})

	if err != nil {
		logrus.Println(err, "it's in parse with claims")
		return 0, nil, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		userId = claims.userId
		userRole = claims.role
		return userId, &userRole, nil
	}

	return 0, nil, err
}
