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
	api.PUT("/categories/:id", controllers.UpdateCategory)

	// Products
	api.GET("/products", controllers.ListProducts)
	api.POST("/products", controllers.CreateProduct)
	api.GET("/products/:id", controllers.GetProduct)
	api.PUT("/products/:id", controllers.UpdateProduct)
	api.DELETE("/products/:id", controllers.DeleteProduct)
	api.POST("/products/:id/image", controllers.UploadProductImage) // New image upload route

	// Cart
	cart := api.Group("/cart")
	{
		cart.POST("", controllers.AddToCart)            // Add item
		cart.GET("", controllers.GetCart)               // List items
		cart.PUT("/:id", controllers.UpdateCartItem)    // Update quantity
		cart.DELETE("/:id", controllers.DeleteCartItem) // Delete one
		cart.DELETE("", controllers.ClearCart)          // Empty cart
	}

	// Order
	order := api.Group("/order")
	{
		order.POST("", controllers.CreateOrder)   // Create from cart
		order.GET("", controllers.GetOrder)       // Latest order
		order.GET("/all", controllers.ListOrders) // (admin) all orders
		order.PUT("", controllers.UpdateOrder)    // Mark paid/unpaid
		order.DELETE("", controllers.DeleteOrder) // Delete all for session
	}
}

