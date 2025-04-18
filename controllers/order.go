package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"bogbon-api/models"
	"bogbon-api/repository"
	"bogbon-api/utils"
)

// CreateOrder is only called manually after cart is filled
func CreateOrder(c *gin.Context) {
	sessionID := utils.GetSessionID(c)

	items, err := repository.GetCartItems(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if len(items) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "cart is empty"})
		return
	}

	order := models.Order{
		SessionID: sessionID,
		IsPaid:    false,
		Items:     items,
	}
	if err := repository.CreateOrder(&order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

func GetOrder(c *gin.Context) {
	sessionID := utils.GetSessionID(c)

	order, err := repository.GetOrderBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// List all orders (for admin or dashboard maybe)
func ListOrders(c *gin.Context) {
	orders, err := repository.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// Update order (e.g., mark as paid)
func UpdateOrder(c *gin.Context) {
	var input struct {
		IsPaid bool `json:"is_paid"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sessionID := utils.GetSessionID(c)
	order, err := repository.GetOrderBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	order.IsPaid = input.IsPaid
	if err := repository.UpdateOrder(order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, order)
}

// Delete order (by session)
func DeleteOrder(c *gin.Context) {
	sessionID := utils.GetSessionID(c)

	if err := repository.DeleteOrderBySession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "order deleted"})
}

