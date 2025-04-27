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
// @Param input body requests.AddToCartInput true "Product to add"
// @Success      201    {object}  models.CartItem
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /cart [post]
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

// GetCart godoc
// @Summary      Get cart items
// @Description  Returns all items in the user's cart.
// @Tags         Cart
// @Produce      json
// @Success      200 {array} models.CartItem
// @Failure      500 {object} map[string]string
// @Router       /cart [get]
func GetCart(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	cart, err := repository.GetCart(sessionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cart.Items)
}

// UpdateCartItem godoc
// @Summary      Update cart item
// @Description  Updates the quantity of a specific cart item.
// @Tags         Cart
// @Accept       json
// @Produce      json
// @Param        id     path      int  true  "Cart Item ID"
// @Param input body requests.UpdateCartItemInput true "Updated quantity"
// @Success      200    {object}  map[string]string
// @Failure      400    {object}  map[string]string
// @Failure      500    {object}  map[string]string
// @Router       /cart/{id} [put]
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

// DeleteCartItem godoc
// @Summary      Delete cart item
// @Description  Deletes a specific cart item from the cart.
// @Tags         Cart
// @Produce      json
// @Param        id   path      int  true  "Cart Item ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Router       /cart/{id} [delete]
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

// ClearCart godoc
// @Summary      Clear cart
// @Description  Deletes all items in the user's cart.
// @Tags         Cart
// @Produce      json
// @Success      204 {object} nil
// @Failure      500 {object} map[string]string
// @Router       /cart/clear [delete]
func ClearCart(c *gin.Context) {
	sessionID := utils.GetSessionID(c)
	if err := repository.ClearCart(sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

