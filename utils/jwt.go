package utils

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"epi/models" 
	"errors"
)

// Secret key for signing JWT tokens (should be stored securely)
var jwtSecret = []byte("your-secret-key-here")

// GenerateToken generates a new JWT token for the user
func GenerateToken(user models.User) (string, error) {
	claims := jwt.MapClaims{
		"sub":    user.Username, // Use Username as the unique identifier
		"role":   user.Role,
		"exp":    time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

// ParseToken parses and validates the JWT token
func ParseToken(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Check if the signing method matches
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	// Token is valid, extract the claims
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
