package services

import (
	"errors"
	"fmt"
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

func (s *AuthService) VerifyAccessToken(tokenString string) (int, error) {
	var userId int
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			logrus.Printf("%s\n", errors.New("unexpected signing method"))
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.signingKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			logrus.Printf("%s\n", errors.New("access token expired"))
			return 0, errors.New("access token expired")
		}
		return 0, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		userId = claims.userId
		logrus.Printf("%s\n", err)
		return 0, err
	}

	return userId, nil
}
