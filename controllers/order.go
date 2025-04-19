package controllers

import (
	"net/http"

	"bogbon-api/repository"
	"bogbon-api/utils"

	"github.com/gin-gonic/gin"
)

// CreateOrder godoc
// @Summary Create a new order from the current cart
// @Tags Orders
// @Produce json
// @Success 201 {object} models.Order
// @Failure 400 {object} map[string]string
// @Router /orders [post]

// CreateOrder creates an order from the current cart.
func CreateOrder(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	order, err := repository.CreateOrderFromCart(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, order)
}

// GetOrder godoc
// @Summary Get the most recent order for the current session
// @Tags Orders
// @Produce json
// @Success 200 {object} models.Order
// @Failure 404 {object} map[string]string
// @Router /orders [get]

// GetOrder returns the most recent order for this session.
func GetOrder(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	order, err := repository.GetOrderBySession(sessionID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}
	c.JSON(http.StatusOK, order)
}

// ListOrders godoc
// @Summary List all orders (admin use)
// @Tags Orders
// @Produce json
// @Success 200 {array} models.Order
// @Failure 500 {object} map[string]string
// @Router /orders/all [get]

// ListOrders returns all orders (admin use).
func ListOrders(c *gin.Context) {
	orders, err := repository.GetAllOrders()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)
}

// UpdateOrder godoc
// @Summary Update the payment status of an order
// @Tags Orders
// @Accept json
// @Produce json
// @Param status body struct{IsPaid bool `json:"is_paid"`} true "Order payment status"
// @Success 200 {object} models.Order
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [put]

// UpdateOrder marks an order paid/unpaid.
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

// DeleteOrder godoc
// @Summary Delete all orders for the current session
// @Tags Orders
// @Produce json
// @Success 200 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /orders [delete]

// DeleteOrder deletes all orders for this session.
func DeleteOrder(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	if err := repository.DeleteOrderBySession(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "orders deleted"})
}

