package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"portal/controllers"
)

func AuthMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}
		claims, err := controllers.ValidateToken(tokenString)
		if err != nil || claims.Role != role {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		c.Set("username", claims.Username)
		c.Next()
	}
}
