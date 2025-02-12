package services

import (
	"context"
	"errors"
	"time"

	"epi/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type AuthService struct {
	DB         *mongo.Client
	Collection *mongo.Collection
}

func NewAuthService(db *mongo.Client) *AuthService {
	collection := db.Database("epi").Collection("users")
	return &AuthService{
		DB:         db,
		Collection: collection,
	}
}

// SaveUser inserts a new user into the database.
func (as *AuthService) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if a user with the same phone number already exists.
	var existing models.User
	err := as.Collection.FindOne(ctx, bson.M{"phone_number": user.PhoneNumber}).Decode(&existing)
	if err == nil {
		return errors.New("user already exists")
	}

	_, err = as.Collection.InsertOne(ctx, user)
	return err
}

// GetUserByPhoneNumber retrieves a user by phone number.
func (as *AuthService) GetUserByPhoneNumber(phoneNumber string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := as.Collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// SaveOtpCode updates the OTPCode field for the given phone number.
func (as *AuthService) SaveOtpCode(phoneNumber, otp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"otp_code":   otp,
			"updated_at": time.Now(),
		},
	}
	res, err := as.Collection.UpdateOne(ctx, bson.M{"phone_number": phoneNumber}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// VerifyOtpCode checks if the provided OTP matches the stored OTP.
func (as *AuthService) VerifyOtpCode(phoneNumber, otp string) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := as.Collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)
	if err != nil {
		return false, errors.New("user not found")
	}
	if user.OTPCode != otp {
		return false, nil
	}
	return true, nil
}

// UpdateUserRole updates the user's role.
func (as *AuthService) UpdateUserRole(phoneNumber, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"role":       role,
			"updated_at": time.Now(),
		},
	}
	res, err := as.Collection.UpdateOne(ctx, bson.M{"phone_number": phoneNumber}, update)
	if err != nil {
		return err
	}
	if res.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// LoginUser verifies the user's password and returns the user record.
func (as *AuthService) LoginUser(phoneNumber, password string) (*models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user models.User
	err := as.Collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)
	if err != nil {
		return nil, errors.New("user not found")
	}

	// In production, compare hashed passwords.
	if user.Password != password {
		return nil, errors.New("incorrect password")
	}
	return &user, nil
}
