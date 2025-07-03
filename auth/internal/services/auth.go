package services

import (
	"errors"
	"fmt"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	uuid "github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	TokenTTL = time.Minute * 60
)

type tokenClaims struct {
	userId   int
	username string
	role     string
	jwt.RegisteredClaims
}

type AuthService struct {
	repo       repository.Authorization
	signingKey string
}

func NewAuthService(repo repository.Authorization, signingKey string) *AuthService {
	return &AuthService{
		repo:       repo,
		signingKey: signingKey,
	}
}

func generatePasswordHash(password string) (string, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(passwordHash), nil
}

func (s *AuthService) Create(user models.User) (int, error) {
	var err error
	user.Password, err = generatePasswordHash(user.Password)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(user)
}

func (s *AuthService) GenerateTokens(username, password string) (string, string, error) {
	user, err := s.repo.GetUser(username)
	if err != nil {
		return "", "", errors.New("internal error")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	refreshToken, err := s.generateRefreshToken(user.Id)
	if err != nil {
		return "", "", err
	}

	accessToken, err := s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil

}

func (s *AuthService) VerifyAccessToken(tokenString string) (models.User, error) {
	user := models.User{}
	token, err := jwt.ParseWithClaims(tokenString, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(s.signingKey), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return models.User{}, errors.New("access token expired")
		}
		return models.User{}, err
	}

	if claims, ok := token.Claims.(*tokenClaims); ok && token.Valid {
		user.Id = claims.userId
		user.Username = claims.username
		user.Role = claims.role
		return user, nil
	}

	return models.User{}, err
}

func (s *AuthService) RefreshTokens(refreshToken string) (string, string, error) {
	tokenRecord, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil || tokenRecord.Revoked || tokenRecord.ExpiresAt.Before(time.Now()) {
		logrus.Println(err)
		return "", "", err
	}

	err = s.repo.RevokeRefreshToken(refreshToken)
	if err != nil {
		logrus.Println(err, "RevokeRefreshToken")
		return "", "", errors.New("can't revoke token")
	}

	user, err := s.repo.GetUserById(tokenRecord.UserId)
	if err != nil {
		logrus.Println(err, "GetUserById")
		return "", "", errors.New("internal error")
	}

	newRefresh, err := s.generateRefreshToken(user.Id)
	if err != nil {
		logrus.Println(err, "generateRefreshToken")
		return "", "", err
	}

	newAccess, err := s.generateAccessToken(user)
	if err != nil {
		logrus.Println(err, "generateAccessToken")
		return "", "", err
	}
	return newAccess, newRefresh, nil
}

func (s *AuthService) generateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		userId:   user.Id,
		username: user.Username,
		role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})
	return token.SignedString([]byte(s.signingKey))
}

func (s *AuthService) generateRefreshToken(userId int) (string, error) {
	token := uuid.New().String()
	expiresAt := time.Now().Add(time.Hour * 7 * 24)

	err := s.repo.SaveRefreshToken(models.RefreshToken{
		UserId:    userId,
		Token:     token,
		ExpiresAt: expiresAt,
	})
	if err != nil {
		return "", err
	}

	return token, nil
}
