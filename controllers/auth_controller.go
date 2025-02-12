package controllers

import (
	"net/http"
	"epi/models"
	"epi/utils"
	"epi/services"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

// SignUp allows customers to sign up with phone number and receive an OTP
func SignUp(c *gin.Context) {
	var signUpData models.User
	if err := c.ShouldBindJSON(&signUpData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Generate OTP and send it via SMS (simulated)
	otp, err := utils.SendOTP(signUpData.PhoneNumber)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send OTP"})
		return
	}

	// Store the OTP temporarily (in real apps, you would store this in a cache/db)
	// For now, assume OTP is valid for 5 minutes
	c.Set("otp", otp)

	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully. Please verify your phone number."})
}

// VerifyOTP allows users to verify the OTP
func VerifyOTP(c *gin.Context) {
	var otpData struct {
		PhoneNumber string `json:"phone_number"`
		OTP         string `json:"otp"`
	}

	if err := c.ShouldBindJSON(&otpData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Retrieve OTP from context (this is a dummy implementation, replace with real logic)
	storedOtp, _ := c.Get("otp")

	// Check if OTP matches the one sent (dummy logic)
	if storedOtp == otpData.OTP {
		// Save the user to the database (this should actually happen in a real service)
		userService := services.NewUserService()
		user := models.User{
			Username:    otpData.PhoneNumber, // Treat phone number as the username
			Password:    "",                  // Password should be handled separately
			PhoneNumber: otpData.PhoneNumber,
			Role:        "customer",
		}
		userService.CreateUser(user)

		// Return success message
		c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully. You can now log in."})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
	}
}

// Login allows customers to log in with phone number and OTP
func Login(c *gin.Context) {
	var loginData struct {
		PhoneNumber string `json:"phone_number"`
		OTP         string `json:"otp"`
	}
	if err := c.ShouldBindJSON(&loginData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Verify OTP
	storedOtp, _ := c.Get("otp")
	if storedOtp == loginData.OTP {
		// Fetch user by phone number
		userService := services.NewUserService()
		user, err := userService.GetUserByPhoneNumber(loginData.PhoneNumber)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid phone number or OTP"})
			return
		}

		// Generate JWT token
		token, err := utils.GenerateToken(user)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
			return
		}

		// Return the token
		c.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
	}
}
