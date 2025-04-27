package requests

// AddToCartInput is the request body for adding an item to cart
type AddToCartInput struct {
	ProductID uint `json:"product_id" binding:"required"`
	Quantity  int  `json:"quantity" binding:"gte=1"`
}

// UpdateCartItemInput is the request body for updating cart item quantity
type UpdateCartItemInput struct {
	Quantity int `json:"quantity" binding:"gte=1"`
}

