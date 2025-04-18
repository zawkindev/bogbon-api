package controllers

import (
	"net/http"
	"strconv"

	"bogbon-api/models"
	"bogbon-api/repository"
	"bogbon-api/utils"
	"github.com/gin-gonic/gin"
)

// Add item to cart
func AddToCart(c *gin.Context) {
	var input struct {
		ProductID uint `json:"product_id" binding:"required"`
		Quantity  int  `json:"quantity" binding:"gte=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := utils.GetSessionID(c)

	// Ensure order exists for foreign key constraint
	if _, err := repository.EnsureOrder(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order"})
		return
	}

	item := models.CartItem{
		SessionID: sessionID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}

	if err := repository.AddCartItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// Get all cart items
func GetCart(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	items, err := repository.GetCartItems(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, items)
}

// Update quantity of a cart item
func UpdateCartItem(c *gin.Context) {
	id := c.Param("id")

	var input struct {
		Quantity int `json:"quantity" binding:"gte=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	itemID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	if err := repository.UpdateCartItem(uint(itemID), input.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusOK)
}

// Delete a specific item from cart
func DeleteCartItem(c *gin.Context) {
	id := c.Param("id")

	itemID, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	if err := repository.DeleteCartItem(uint(itemID)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// Clear entire cart for a session
func ClearCart(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	if err := repository.ClearCart(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

