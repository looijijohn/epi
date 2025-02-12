package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "epi/models"
    "epi/services"
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