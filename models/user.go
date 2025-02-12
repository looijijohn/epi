package models

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

 

type User struct {
	ID             bson.ObjectID `bson:"_id,omitempty"`
	PhoneNumber string            `bson:"phone_number"`
	Password   string            `bson:"password,omitempty"`
	Role       string            `bson:"role"`       // Role: "admin" or "customer"
	OTPCode    string            `bson:"otp_code"`   // OTP code for verification
	CreatedAt  time.Time         `bson:"created_at"`
	UpdatedAt  time.Time         `bson:"updated_at"`
}

type SignUpData struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	Password    string `json:"password" binding:"required"`
}

type OtpVerification struct {
	PhoneNumber string `json:"phone_number" binding:"required"`
	OtpCode     string `json:"otp_code" binding:"required"`
}
