package controllers

import (
	"net/http"
	"strconv"

	"bogbon-api/models"
	"bogbon-api/repository"
	"bogbon-api/utils"

	"github.com/gin-gonic/gin"
)

// AddToCart godoc
// @Summary      Add item to cart
// @Description  Adds a product to the user's session cart.
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        input  body      struct{ProductID uint ` + "`json:\"product_id\" binding:\"required\"`" + `;Quantity int ` + "`json:\"quantity\" binding:\"gte=1\"`" + `}  true  "Product to add"
// @Success      201    {object}  models.CartItem
// @Failure      400    {object}  gin.H{"error": "…"}  
// @Failure      500    {object}  gin.H{"error": "…"}  
// @Router       /cart [post]

// AddToCart adds a product to the user's cart.
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
	cart, err := repository.EnsureCart(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not get or create cart"})
		return
	}

	item := models.CartItem{
		CartID:    cart.ID,
		ProductID: input.ProductID,
		Quantity:  input.Quantity,
	}
	if err := repository.AddCartItem(&item); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, item)
}

// GetCart returns all items in the user's cart.
func GetCart(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	cart, err := repository.GetCart(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart.Items)
}

// UpdateCartItem updates the quantity of a cart item.
func UpdateCartItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	var input struct {
		Quantity int `json:"quantity" binding:"gte=1"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := repository.UpdateCartItem(uint(id), input.Quantity); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "cart item updated"})
}

// DeleteCartItem removes a single item from the cart.
func DeleteCartItem(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid item ID"})
		return
	}

	if err := repository.DeleteCartItem(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// ClearCart deletes all items in the cart.
func ClearCart(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	if err := repository.ClearCart(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

