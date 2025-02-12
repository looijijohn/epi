package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "epi/models"
    "epi/services"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type ProductController struct {
    ProductService *services.ProductService
}

func NewProductController(productService *services.ProductService) *ProductController {
    return &ProductController{ProductService: productService}
}

// GetProducts retrieves all products
func (pc *ProductController) GetProducts(c *gin.Context) {
    products, err := pc.ProductService.GetAllProducts()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, products)
}

// CreateProduct creates a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    insertedID, err := pc.ProductService.CreateProduct(product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{"inserted_id": insertedID})
}

// DeleteProduct deletes a product by ID
func (pc *ProductController) DeleteProduct(c *gin.Context) {
    productID := c.Param("id")
    objectID, err := bson.ObjectIDFromHex(productID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    err = pc.ProductService.DeleteProduct(objectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product deleted successfully"})
}

// UpdateProduct updates a product by ID
func (pc *ProductController) UpdateProduct(c *gin.Context) {
    productID := c.Param("id")
    objectID, err := bson.ObjectIDFromHex(productID)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
        return
    }

    var updatedProduct models.Product
    if err := c.ShouldBindJSON(&updatedProduct); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    err = pc.ProductService.UpdateProduct(objectID, updatedProduct)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Product updated successfully"})
}