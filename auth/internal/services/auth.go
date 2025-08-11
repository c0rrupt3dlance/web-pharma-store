package services

import (
	"errors"
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
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})

	if err != nil {
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
	var (
		newRefreshToken string
		newAccessToken  string
	)

	tokenRecord, err := s.repo.GetRefreshToken(refreshToken)
	if err != nil || tokenRecord.Revoked || tokenRecord.ExpiresAt.Before(time.Now()) {
		return "", "", err
	}

	err = s.repo.RevokeRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	user, err := s.repo.GetUserById(tokenRecord.UserId)
	if err != nil {
		return "", "", err
	}

	newRefreshToken, err = s.generateRefreshToken(tokenRecord.UserId)
	if err != nil {
		return "", "", err
	}

	newAccessToken, err = s.generateAccessToken(user)
	if err != nil {
		return "", "", err
	}

	return newRefreshToken, newAccessToken, nil

}

func (s *AuthService) generateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tokenClaims{
		userId:   user.Id,
		username: user.Name,
		role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authService",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s *AuthService) generateRefreshToken(id int) (string, error) {
	str := uuid.New().String()

	var token = models.RefreshToken{
		UserId:    id,
		Token:     str,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
		IssuedAt:  time.Now(),
	}

	err := s.repo.SaveRefreshToken(token)
	if err != nil {
		logrus.Println(err)
	}

	return token.Token, nil
}
