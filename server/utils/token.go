package utils

import (
	"errors"
	"os"
	"strings"
	"time"

	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/models"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken tạo JWT Access Token
func GenerateToken(user models.User) (string, error) {

	claims := JWTClaims{
		UserID: user.ID.Hex(),
		Email:  user.Email,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", errors.New("JWT_SECRET is not configured")
	}

	return token.SignedString([]byte(secret))
}

// ValidateToken kiểm tra và parse JWT
func ValidateToken(tokenString string) (*JWTClaims, error) {

	claims := &JWTClaims{}

	secret := os.Getenv("JWT_SECRET")

	token, err := jwt.ParseWithClaims(
		tokenString,
		claims,
		func(token *jwt.Token) (interface{}, error) {
			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if time.Now().After(claims.ExpiresAt.Time) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

func GetAccessToken(c *gin.Context) (string, error) {

	authHeader := c.GetHeader("Authorization")

	if authHeader == "" {
		return "", errors.New("authorization header is required")
	}

	const bearerPrefix = "Bearer "

	if !strings.HasPrefix(authHeader, bearerPrefix) {
		return "", errors.New("invalid authorization header format")
	}

	token := strings.TrimPrefix(authHeader, bearerPrefix)

	if token == "" {
		return "", errors.New("access token is empty")
	}

	return token, nil
}
