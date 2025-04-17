// internal/controllers/order.go
package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "bogbon-api/models"
    "bogbon-api/repository"
    "bogbon-api/utils"
)

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

