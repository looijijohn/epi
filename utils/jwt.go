package utils

import (
	"errors"
	"time"

	"github.com/dgrijalva/jwt-go"
	"epi/models"
)

// JWT_SECRET is used to sign tokens.
var JWT_SECRET = []byte("your-secret-key")

// JwtCustomClaims defines the custom JWT claims.
type JwtCustomClaims struct {
	ID   string `json:"id"`   // Using the phone number as ID
	Role string `json:"role"` // e.g., "customer" or "admin"
	jwt.StandardClaims
}

// GenerateToken creates a JWT for the given user.
func GenerateToken(user models.User) (string, error) {
	claims := &JwtCustomClaims{
		ID:   user.PhoneNumber,
		Role: user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(1 * time.Hour).Unix(),
			IssuedAt:  time.Now().Unix(),
			Issuer:    "epi",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(JWT_SECRET)
}

// ValidateToken parses and validates the JWT token.
func ValidateToken(tokenStr string) (*JwtCustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JwtCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return JWT_SECRET, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*JwtCustomClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, errors.New("invalid token")
}
