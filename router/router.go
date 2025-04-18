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

	cartGroup := r.Group("/api/cart")
	{
		cartGroup.POST("/", controllers.AddToCart)           // Add an item to the cart
		cartGroup.GET("/", controllers.GetCart)              // Get all cart items for the session
		cartGroup.PUT("/:id", controllers.UpdateCartItem)    // Update cart item (quantity)
		cartGroup.DELETE("/:id", controllers.DeleteCartItem) // Delete a cart item
		cartGroup.DELETE("/", controllers.ClearCart)         // Clear all items from the cart
	}

	order := r.Group("/api/order")
	{
		order.POST("", controllers.CreateOrder)
		order.GET("", controllers.GetOrder)
		order.GET("/all", controllers.ListOrders) // optional
		order.PUT("", controllers.UpdateOrder)
		order.DELETE("", controllers.DeleteOrder)
	}

}
