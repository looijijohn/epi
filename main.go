package main

import (
    "epi/config"
    "epi/controllers"
    "epi/routes"
    "epi/services"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Connect to MongoDB
    db := config.ConnectDB()

    // Initialize services
    productService := services.NewProductService(db)

    // Initialize controllers
    productController := controllers.NewProductController(productService)

    // Register routes
    routes.ProductRoutes(r, productController)

    // Start the server
    r.Run(":8080")
}

