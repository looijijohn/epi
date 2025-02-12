package services

import (
	"context"
	"time"

	"epi/models"
 
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo" 

)

type ProductService struct {
    DB *mongo.Client
}

func NewProductService(db *mongo.Client) *ProductService {
    return &ProductService{DB: db}
}

// GetAllProducts retrieves all products
func (ps *ProductService) GetAllProducts() ([]models.Product, error) {
    collection := ps.DB.Database("epi").Collection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var products []models.Product
    cursor, err := collection.Find(ctx, bson.M{})
    if err != nil {
        return nil, err
    }
    defer cursor.Close(ctx)

    if err = cursor.All(ctx, &products); err != nil {
        return nil, err
    }

    return products, nil
}

// CreateProduct creates a new product
func (ps *ProductService) CreateProduct(product models.Product) ( bson.ObjectID, error) {
    collection := ps.DB.Database("epi").Collection("products")
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    product.CreatedAt = time.Now()
    product.UpdatedAt = time.Now()

    result, err := collection.InsertOne(ctx, product)
    if err != nil {
        return bson.NilObjectID, err
    }

    return result.InsertedID.(bson.ObjectID), nil
}