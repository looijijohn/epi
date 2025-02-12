package controllers

import (
	"net/http"
	"time"

	"epi/models"
	"epi/services"
	"epi/utils"

	"github.com/gin-gonic/gin"
)

// AuthController handles authentication requests.
type AuthController struct {
	AuthService *services.AuthService
}

// NewAuthController creates a new AuthController.
func NewAuthController(authService *services.AuthService) *AuthController {
	return &AuthController{AuthService: authService}
}

// SignupOTP handles sign-up via OTP (phone number only, no password).
func (ac *AuthController) SignupOTP(c *gin.Context) {
	var req struct {
		PhoneNumber string `json:"phone_number" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user exists. If not, create a new one with role "pending".
	user, _ := ac.AuthService.GetUserByPhoneNumber(req.PhoneNumber)
	if user == nil {
		newUser := models.User{
			PhoneNumber: req.PhoneNumber,
			Role:        "pending",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}
		if err := ac.AuthService.SaveUser(newUser); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Generate OTP (dummy)
	otp := utils.GenerateOTP()

	// Save OTP for the user.
	if err := ac.AuthService.SaveOtpCode(req.PhoneNumber, otp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save OTP"})
		return
	}

	// (Optionally) Send the OTP via SMS using utils.SendSMS.
	c.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

// VerifyOTP verifies the OTP sent to the user and returns a JWT token.
func (ac *AuthController) VerifyOTP(c *gin.Context) {
	var req models.OtpVerification
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Verify the OTP.
	valid, err := ac.AuthService.VerifyOtpCode(req.PhoneNumber, req.OtpCode)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired OTP"})
		return
	}

	// Update user's role to "customer".
	if err := ac.AuthService.UpdateUserRole(req.PhoneNumber, "customer"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user role"})
		return
	}

	// Retrieve the user.
	user, err := ac.AuthService.GetUserByPhoneNumber(req.PhoneNumber)
	if err != nil || user == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not found"})
		return
	}

	// Generate JWT token.
	token, err := utils.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "OTP verified successfully", "token": token})
}

// SignupWithPassword registers a new user with a phone number and password.
func (ac *AuthController) SignupWithPassword(c *gin.Context) {
	var req models.SignUpData
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if user already exists.
	user, _ := ac.AuthService.GetUserByPhoneNumber(req.PhoneNumber)
	if user != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	// Create a new user with the password (password should be hashed in production).
	newUser := models.User{
		PhoneNumber: req.PhoneNumber,
		Password:    req.Password,
		Role:        "customer",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := ac.AuthService.SaveUser(newUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token.
	token, err := utils.GenerateToken(newUser)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User registered successfully", "token": token})
}

// LoginWithPassword logs in a user using phone number and password.
func (ac *AuthController) LoginWithPassword(c *gin.Context) {
	var req models.SignUpData // Reusing SignUpData struct: phone_number and password.
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := ac.AuthService.LoginUser(req.PhoneNumber, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Generate JWT token.
	token, err := utils.GenerateToken(*user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "token": token})
}
