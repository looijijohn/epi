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
// Returns an error if a user with the same phone number already exists.
func (as *AuthService) SaveUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if a user with the given phone number already exists.
	var existing models.User
	err := as.Collection.FindOne(ctx, bson.M{"phone_number": user.PhoneNumber}).Decode(&existing)
	if err == nil {
		return errors.New("user already exists")
	}

	// Insert the new user.
	_, err = as.Collection.InsertOne(ctx, user)
	return err
}

// GetUserByPhoneNumber returns the user with the specified phone number.
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

// SaveOtpCode updates the OTPCode field of the user with the given phone number.
func (as *AuthService) SaveOtpCode(phoneNumber, otp string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"otp_code":   otp,
			"updated_at": time.Now(),
		},
	}
	result, err := as.Collection.UpdateOne(ctx, bson.M{"phone_number": phoneNumber}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}

// VerifyOtpCode compares the provided OTP with the stored OTP.
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

// UpdateUserRole updates the role of the user identified by phone number.
func (as *AuthService) UpdateUserRole(phoneNumber, role string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"role":       role,
			"updated_at": time.Now(),
		},
	}
	result, err := as.Collection.UpdateOne(ctx, bson.M{"phone_number": phoneNumber}, update)
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("user not found")
	}
	return nil
}
