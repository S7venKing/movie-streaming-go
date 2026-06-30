package middleware

import (
	"net/http"

	"github.com/S7venKing/movie-streaming-go/server/magic-stream-movie-server/utils"
	"github.com/gin-gonic/gin"
)

func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := utils.GetAccessToken(c)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
		}
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
		}
		c.Set("userId", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
