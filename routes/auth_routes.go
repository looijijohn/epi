package routes

import (
	"github.com/gin-gonic/gin"
	"epi/controllers"
	"epi/services"
)

// AuthRoutes registers auth endpoints.
func AuthRoutes(router *gin.Engine, authService *services.AuthService) {
	authController := controllers.NewAuthController(authService)
	auth := router.Group("/auth")
	{
		auth.POST("/signup-otp", authController.SignupOTP)           // Sign up with OTP (no password)
		auth.POST("/verify-otp", authController.VerifyOTP)             // Verify OTP to get token
		auth.POST("/signup", authController.SignupWithPassword)        // Sign up with phone number & password
		auth.POST("/login", authController.LoginWithPassword)          // Log in with phone number & password
	}
}
