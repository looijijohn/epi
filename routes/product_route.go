package routes

import (
    "github.com/gin-gonic/gin"
    "epi/controllers"
)

func ProductRoutes(router *gin.Engine, productController *controllers.ProductController) {
    productRoutes := router.Group("/products")
    {
        productRoutes.GET("/", productController.GetProducts)
        productRoutes.POST("/", productController.CreateProduct)
    }
}