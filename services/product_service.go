package services

import (
	"context"
	"time"

	"epi/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)


type ProductService struct {
    DB         *mongo.Client
    Collection *mongo.Collection // Add collection as a field
}

func NewProductService(db *mongo.Client) *ProductService {
    collection := db.Database("epi").Collection("products") // Initialize collection once
    return &ProductService{
        DB:         db,
        Collection: collection, // Assign collection to the struct field
    }
}

// GetAllProducts retrieves all products
func (ps *ProductService) GetAllProducts() ([]models.Product, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var products []models.Product
    cursor, err := ps.Collection.Find(ctx, bson.M{}) // Use ps.Collection
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
func (ps *ProductService) CreateProduct(product models.Product) (bson.ObjectID, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    product.CreatedAt = time.Now()
    product.UpdatedAt = time.Now()

    result, err := ps.Collection.InsertOne(ctx, product) // Use ps.Collection
    if err != nil {
        return bson.NilObjectID, err
    }

    return result.InsertedID.(bson.ObjectID), nil
}

// DeleteProduct deletes a product by ID
func (ps *ProductService) DeleteProduct(productID bson.ObjectID) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"_id": productID}
    _, err := ps.Collection.DeleteOne(ctx, filter) // Use ps.Collection
    if err != nil {
        return err
    }

    return nil
}

// UpdateProduct updates a product by ID
func (ps *ProductService) UpdateProduct(productID bson.ObjectID, updatedProduct models.Product) error {
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    filter := bson.M{"_id": productID}
    update := bson.M{
        "$set": bson.M{
            "name":        updatedProduct.Name,
            "description": updatedProduct.Description,
            "price":       updatedProduct.Price,
            "category_id": updatedProduct.CategoryID,
            "stock":       updatedProduct.Stock,
            "updated_at":  time.Now(),
        },
    }

    _, err := ps.Collection.UpdateOne(ctx, filter, update) // Use ps.Collection
    if err != nil {
        return err
    }

    return nil
}