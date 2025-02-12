package main

import (
    "epi/config"
    "epi/controllers"
    "epi/routes"
    "epi/services"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()

    // Connect to MongoDB
    dbCli := config.ConnectDB()


    
    // Initialize services
    productService := services.NewProductService(dbCli)
    // Initialize controllers
    productController := controllers.NewProductController(productService)
    // Register routes
    routes.ProductRoutes(router, productController)




    // Initialize services
    authService  := services.NewAuthService(dbCli)
	// Register routes
	routes.AuthRoutes(router, authService) // Auth routes


    // Start the server
    router.Run(":8080")
}

