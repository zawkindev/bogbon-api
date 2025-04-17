package router

import (
	"bogbon-api/controllers"

	"github.com/gin-gonic/gin"
)

func Setup(r *gin.Engine) {
	api := r.Group("/api")

	// Categories
	api.GET("/categories", controllers.ListCategories)
	api.POST("/categories", controllers.CreateCategory)
	api.DELETE("/categories/:id", controllers.DeleteCategory)

	// Products
	api.GET("/products", controllers.ListProducts)
	api.POST("/products", controllers.CreateProduct)
	api.GET("/products/:id", controllers.GetProduct)
	api.PUT("/products/:id", controllers.UpdateProduct)
	api.DELETE("/products/:id", controllers.DeleteProduct)

	// Cart & Orders
	api.POST("/cart", controllers.AddToCart)
	api.GET("/cart", controllers.GetCart)
	api.POST("/order", controllers.CreateOrder)
	api.GET("/order", controllers.GetOrder)
}
