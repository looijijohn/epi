package services

import (
	"epi/models"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/bson"
	"context"
	"time"
)

type UserService struct {
    DB         *mongo.Client
    Collection *mongo.Collection
}

func NewUserService(db *mongo.Client) *UserService {
    collection := db.Database("epi").Collection("users")
    return &UserService{
        DB:         db,
        Collection: collection,
    }
}

// GetUserByPhoneNumber fetches a user by their phone number
func (us *UserService) GetUserByPhoneNumber(phoneNumber string) (models.User, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var user models.User
    err := us.Collection.FindOne(ctx, bson.M{"phone_number": phoneNumber}).Decode(&user)
    if err != nil {
        return models.User{}, err
    }
    return user, nil
}

// CreateUser creates a new user in the database
func (us *UserService) CreateUser(user models.User) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    _, err := us.Collection.InsertOne(ctx, user)
    if err != nil {
        return err
    }
    return nil
}
