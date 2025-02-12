package routes

import (
	"github.com/gin-gonic/gin"
	"epi/controllers"
	"epi/services"
)

// AuthRoutes registers the authentication endpoints.
func AuthRoutes(router *gin.Engine, authService *services.AuthService) {
	authController := controllers.NewAuthController(authService)

	auth := router.Group("/auth")
	{
		auth.POST("/signup", authController.SignUp)
		auth.POST("/verify-otp", authController.VerifyOTP)
		auth.POST("/generate-otp", authController.GenerateOTP) // Optional endpoint to re-generate OTP.
	}
}
