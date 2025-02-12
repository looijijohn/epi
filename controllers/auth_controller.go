package controllers

import (
	"net/http"
	"time"

	"epi/models"
	"epi/services"
	"epi/utils"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication-related HTTP requests.
type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController creates a new AuthController.
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// SignUp registers a new user and sends an OTP.
func (ac *AuthController) SignUp(c *gin.Context) {
	var signUpData models.SignUpData
	if err := c.ShouldBindJSON(&signUpData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create a new user with a "pending" role.
	user := models.User{
		PhoneNumber: signUpData.PhoneNumber,
		Password:    signUpData.Password, // In production, hash the password!
		Role:        "pending",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Save the user (fails if the phone number already exists).
	if err := ac.AuthService.SaveUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Generate a dummy OTP.
	otp := utils.GenerateOTP()

	// Save the OTP in the user's record.
	if err := ac.AuthService.SaveOtpCode(signUpData.PhoneNumber, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}

	// (Optionally) send the OTP via SMS using utils.SendSMS.

	c.JSON(http.StatusOK, gin.H{"message": "User registered. OTP sent to phone number."})
}

// VerifyOTP verifies the provided OTP; if valid, it updates the user’s role and returns a JWT.
func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var otpData models.OtpVerification
	if err := c.ShouldBindJSON(&otpData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify that the OTP is correct.
	valid, err := ac.AuthService.VerifyOtpCode(otpData.PhoneNumber, otpData.OtpCode)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Retrieve the user record.
	user, err := ac.AuthService.GetUserByPhoneNumber(otpData.PhoneNumber)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Update the user role to "customer" after successful OTP verification.
	if err := ac.AuthService.UpdateUserRole(otpData.PhoneNumber, "customer"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	// Generate a JWT token.
	token, err := utils.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully", "token": token})
}

// GenerateOTP is an optional endpoint to regenerate an OTP.
func (ac *AuthController) GenerateOTP(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	otp := utils.GenerateOTP()
	if err := ac.AuthService.SaveOtpCode(req.PhoneNumber, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate OTP"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP generated and saved"})
}
