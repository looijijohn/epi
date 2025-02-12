package middleware

import (
	"net/http"
	"strings"
	"epi/utils"
	"github.com/gin-gonic/gin" 
)

// AuthMiddleware is a middleware for JWT token authentication
func AuthMiddleware(requiredRole string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Extract token from Authorization header
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			c.Abort()
			return
		}

		// Remove the "Bearer " prefix from the token string
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")

		// Parse and validate the token
		claims, err := utils.ParseToken(tokenString)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		// Check if the user has the required role
		if requiredRole != "" && claims["role"] != requiredRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "Forbidden: Insufficient role"})
			c.Abort()
			return
		}

		// Pass user claims to the context
		c.Set("userID", claims["sub"])
		c.Set("role", claims["role"])

		// Continue with the request
		c.Next()
	}
}
