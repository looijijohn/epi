package services

import (
	"context"
	"fmt"
	"time"

	"epi/models"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

// ProductService handles product-related database operations.
type ProductService struct {
	DB         *mongo.Client
	Collection *mongo.Collection
}

func NewProductService(db *mongo.Client) *ProductService {
	collection := db.Database("epi").Collection("products")
	return &ProductService{
		DB:         db,
		Collection: collection,
	}
}

// GetAllProducts retrieves all products from the database.
func (ps *ProductService) GetAllProducts() ([]models.Product, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cursor, err := ps.Collection.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var products []models.Product
	if err := cursor.All(ctx, &products); err != nil {
		return nil, err
	}
	return products, nil
}

// CreateProduct inserts a new product into the database.
func (ps *ProductService) CreateProduct(product *models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := ps.Collection.InsertOne(ctx, product)
	return err
}

// UpdateProduct updates an existing product by its ID.
func (ps *ProductService) UpdateProduct(id string, updatedProduct models.Product) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the id string to a BSON ObjectID.
	// Convert the id string to a primitive.ObjectID.
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}


	filter := bson.M{"_id": objID}
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
	_, err = ps.Collection.UpdateOne(ctx, filter, update)
	return err
}

// DeleteProduct removes a product from the database by its ID.
// DeleteProduct removes a product from the database by its ID.
func (ps *ProductService) DeleteProduct(id string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Convert the id string to a BSON ObjectID.
	objID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return err
	}

	filter := bson.M{"_id": objID}
	result, err := ps.Collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	// Check if no documents were deleted
	if result.DeletedCount == 0 {
		return fmt.Errorf("no product found with the given ID")
	}

	return nil
}
