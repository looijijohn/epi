package models

import "go.mongodb.org/mongo-driver/v2/bson"

// User model represents a user in the system (Admin or Customer)
type User struct {
    ID          bson.ObjectID `bson:"_id,omitempty"`
    Username    string        `bson:"username"`
    Password    string        `bson:"password"` // Password will be hashed
    PhoneNumber string        `bson:"phone_number"`
    Role        string        `bson:"role"`     // "admin" or "customer"
}
