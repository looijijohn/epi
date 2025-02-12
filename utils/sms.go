package utils

import (
	"math/rand"
	"strconv"
	"time"
)

// GenerateOTP returns a random 6-digit OTP as a string.
func GenerateOTP() string {
	rand.Seed(time.Now().UnixNano())
	otp := rand.Intn(900000) + 100000 // Generates a number between 100000 and 999999.
	return strconv.Itoa(otp)
}

// Optionally, you can add a dummy SendSMS function.
func SendSMS(phoneNumber, message string) error {
	// For now, simply log or print the message.
	// In production, integrate with an SMS API.
	return nil
}
