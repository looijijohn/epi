package models

import (
    "go.mongodb.org/mongo-driver/v2/bson"
    "time"
)

type Product struct {
    ID          bson.ObjectID `bson:"_id,omitempty"`
    Name        string             `bson:"name"`
    Description string             `bson:"description"`
    Price       float64            `bson:"price"`
    CategoryID  bson.ObjectID `bson:"category_id"`
    Stock       int                `bson:"stock"`
    CreatedAt   time.Time          `bson:"created_at"`
    UpdatedAt   time.Time          `bson:"updated_at"`
}