package middleware

import (
	"net/http"
	"strings"
	"time"

	"epi/utils"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware checks for a valid JWT token.
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Get the Authorization header (should be "Bearer <token>")
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header missing"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header format must be Bearer {token}"})
			c.Abort()
			return
		}

		tokenStr := parts[1]
		claims, err := utils.ValidateToken(tokenStr)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token: " + err.Error()})
			c.Abort()
			return
		}

		// Check if token has expired.
		if time.Now().Unix() > claims.ExpiresAt {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		// Set user information in context for downstream handlers.
		c.Set("user_id", claims.ID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
