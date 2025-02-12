package controllers

import (
    "net/http"
    "epi/models"
    "epi/services"
    "epi/utils" // Import the utils package
    "go.mongodb.org/mongo-driver/v2/bson"
    "github.com/gin-gonic/gin"
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
        c.JSON(http.StatusInternalServerError, utils.ApiErrorResponse{
            Status:  false,
            Message: err.Error(),
        })
        return
    }
    c.JSON(http.StatusOK, utils.ApiResponse{
        Status: true,
        Data:   products,
    })
}

// CreateProduct creates a new product
func (pc *ProductController) CreateProduct(c *gin.Context) {
    var product models.Product
    if err := c.ShouldBindJSON(&product); err != nil {
        c.JSON(http.StatusBadRequest, utils.ApiErrorResponse{
            Status:  false,
            Message: err.Error(),
        })
        return
    }

    insertedID, err := pc.ProductService.CreateProduct(product)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ApiErrorResponse{
            Status:  false,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusCreated, utils.ApiResponse{
        Status: true,
        Data:   gin.H{"inserted_id": insertedID},
    })
}

// DeleteProduct deletes a product by ID
func (pc *ProductController) DeleteProduct(c *gin.Context) {
    productID := c.Param("id")
    objectID, err := bson.ObjectIDFromHex(productID)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ApiErrorResponse{
            Status:  false,
            Message: "Invalid product ID",
        })
        return
    }

    err = pc.ProductService.DeleteProduct(objectID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ApiErrorResponse{
            Status:  false,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, utils.ApiResponse{
        Status: true,
        Data:   gin.H{"message": "Product deleted successfully"},
    })
}

// UpdateProduct updates a product by ID
func (pc *ProductController) UpdateProduct(c *gin.Context) {
    productID := c.Param("id")
    objectID, err := bson.ObjectIDFromHex(productID)
    if err != nil {
        c.JSON(http.StatusBadRequest, utils.ApiErrorResponse{
            Status:  false,
            Message: "Invalid product ID",
        })
        return
    }

    var updatedProduct models.Product
    if err := c.ShouldBindJSON(&updatedProduct); err != nil {
        c.JSON(http.StatusBadRequest, utils.ApiErrorResponse{
            Status:  false,
            Message: err.Error(),
        })
        return
    }

    err = pc.ProductService.UpdateProduct(objectID, updatedProduct)
    if err != nil {
        c.JSON(http.StatusInternalServerError, utils.ApiErrorResponse{
            Status:  false,
            Message: err.Error(),
        })
        return
    }

    c.JSON(http.StatusOK, utils.ApiResponse{
        Status: true,
        Data:   gin.H{"message": "Product updated successfully"},
    })
}
