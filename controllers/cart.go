// internal/controllers/cart.go
package controllers

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "bogbon-api/models"
    "bogbon-api/repository"
    "bogbon-api/utils"
)

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

func GetCart(c *gin.Context) {
    sessionID := utils.GetSessionID(c)
    items, err := repository.GetCartItems(sessionID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, items)
}

