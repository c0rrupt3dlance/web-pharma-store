package services

import (
	"errors"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/models"
	"github.com/c0rrupt3dlance/web-pharma-store/auth/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"time"
)

const (
	TokenTTL = time.Minute * 60
)

type accessTokenClaims struct {
	UserId   int    `json:"user_id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

type refreshTokenClaims struct {
	UserId int `json:"user_id"`
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
	logrus.Println(tokenString)
	token, err := jwt.ParseWithClaims(tokenString, &accessTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
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

	if claims, ok := token.Claims.(*accessTokenClaims); ok && token.Valid {
		user.Id = claims.UserId
		user.Username = claims.Username
		user.Role = claims.Role
		return user, nil
	}

	return models.User{}, err
}

func (s *AuthService) RefreshTokens(refreshTokenString string) (string, string, error) {
	var (
		newRefreshToken string
		newAccessToken  string
	)
	token, err := jwt.ParseWithClaims(refreshTokenString, &refreshTokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)

		if !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.signingKey), nil
	})

	if err != nil {
		return "", "", err
	}

	claims, ok := token.Claims.(*refreshTokenClaims)
	if !ok || !token.Valid {
		return "", "", err
	}

	refreshTokenRecorded, err := s.repo.GetRefreshToken(refreshTokenString)
	if err != nil {
		logrus.Println(err)
		return newRefreshToken, newAccessToken, errors.New("token not valid or not found")
	}

	if refreshTokenRecorded.Revoked || claims.ExpiresAt.Before(time.Now()) {
		err = s.repo.RevokeRefreshToken(refreshTokenRecorded.Token)
		if err != nil {
			logrus.Println(err)
			return newRefreshToken, newAccessToken, errors.New("expired or revoked token")
		}
	}

	err = s.repo.RevokeRefreshToken(refreshTokenRecorded.Token)
	if err != nil {
		logrus.Println(err)
		return newRefreshToken, newAccessToken, errors.New("couldn't refresh tokens")
	}

	user, err := s.repo.GetUserById(refreshTokenRecorded.UserId)
	if err != nil {
		logrus.Println(err)
		return newRefreshToken, newAccessToken, errors.New("couldn't refresh tokens")
	}

	newRefreshToken, err = s.generateRefreshToken(refreshTokenRecorded.UserId)
	if err != nil {
		logrus.Println(err)
		return newRefreshToken, newAccessToken, errors.New("couldn't refresh tokens")
	}

	newAccessToken, err = s.generateAccessToken(user)
	if err != nil {
		logrus.Println(err)
		return newRefreshToken, newAccessToken, errors.New("couldn't refresh tokens")
	}

	return newRefreshToken, newAccessToken, nil
}

func (s *AuthService) generateAccessToken(user models.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenClaims{
		UserId:   user.Id,
		Username: user.Name,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    "authService",
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	})

	return token.SignedString([]byte(s.signingKey))
}

func (s *AuthService) generateRefreshToken(id int) (string, error) {
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.RegisteredClaims{
		Issuer:    "authService",
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenTTL * 24 * 7)),
	})

	str, err := refreshToken.SignedString([]byte(s.signingKey))
	if err != nil {
		return "", errors.New("couldn't get string of a refresh token")
	}

	var token = models.RefreshToken{
		UserId: id,
		Token:  str,
	}

	err = s.repo.SaveRefreshToken(token)
	if err != nil {
		return "", err
	}

	return token.Token, nil
}
