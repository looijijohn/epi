package routes

import (
	"epi/controllers"
	"epi/middleware"

	"github.com/gin-gonic/gin"
)

// ProductRoutes registers product-related endpoints.
func ProductRoutes(router *gin.Engine, productController *controllers.ProductController) {
	productRoutes := router.Group("/products")
	{
		// Public route: list products.
		productRoutes.GET("/", productController.GetProducts)

		// Protected routes (apply your AuthMiddleware to enforce JWT validation).
		// For example, if you have a middleware package, you might add:
		productRoutes.Use(middleware.AuthMiddleware())
		productRoutes.POST("/", productController.CreateProduct)
		productRoutes.PUT("/:id", productController.UpdateProduct)
		productRoutes.DELETE("/:id", productController.DeleteProduct)
	}
}
