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
			return
		}

		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"message": "No token provided"})
			c.Abort()
			return
		}

		claims, err := utils.ValidateToken(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"message": err.Error()})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("role", claims.Role)

		c.Next()
	}
}

func RoleMiddleware(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		role, exists := c.Get("role")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Role not found",
			})
			c.Abort()
			return
		}

		userRole, ok := role.(string)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "Invalid role",
			})
			c.Abort()
			return
		}

		for _, r := range roles {
			if r == userRole {
				c.Next()
				return
			}
		}

		c.JSON(http.StatusForbidden, gin.H{
			"message": "Permission denied",
		})
		c.Abort()
	}
}
