package models

import (
	"gorm.io/gorm"
	"time"
)

// Category: fixed taxonomy
type Category struct {
	ID           uint                  `gorm:"primaryKey;autoIncrement"`
	Translations []CategoryTranslation `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Products     []Product             `gorm:"many2many:category_products;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CategoryTranslation struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	CategoryID   uint   `gorm:"not null"`
	LanguageCode string `gorm:"size:10;not null"` // e.g., "en", "es"
	Name         string `gorm:"not null"`
}

// Product: can be "plant" or "service" determined by Type
type Product struct {
	ID           uint                 `gorm:"primaryKey;autoIncrement"`
	Price        int                  `gorm:"not null"`
	Stock        int                  `gorm:"not null"`
	Type         string               `gorm:"type:VARCHAR(20);not null"`
	Categories   []Category           `gorm:"many2many:category_products;"`
	Translations []ProductTranslation `gorm:"foreignKey:ProductID"`
	Images       []ProductImage       `gorm:"foreignKey:ProductID;constraint:OnDelete:CASCADE;"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`
}

// ProductTranslation: translations for product name and description in different languages
type ProductTranslation struct {
	ID           uint   `gorm:"primaryKey;autoIncrement"`
	ProductID    uint   `gorm:"not null"`
	LanguageCode string `gorm:"size:10;not null"` // For example, "en", "ru", "uz"
	Name         string `gorm:"not null"`
	Description  string
	ShortInfo    string
}

// image
type ProductImage struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	ProductID uint   `gorm:"not null;index"`
	URL       string `gorm:"size:255;not null"`
	CreatedAt time.Time
}

// Cart model: holds the cart items before checkout
type Cart struct {
	ID        uint       `gorm:"primaryKey;autoIncrement"`
	SessionID string     `gorm:"index;not null;unique"`
	Items     []CartItem `gorm:"foreignKey:CartID"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

// CartItem model: now belongs to a Cart (not directly to SessionID or OrderID)
type CartItem struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	CartID    uint `gorm:"index;not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null;default:1"`
	Product   Product
}

// Order model: created from a Cart
type Order struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	SessionID string `gorm:"index;not null"`
	CartID    uint   `gorm:"not null"`
	IsPaid    bool   `gorm:"default:false"`
	CreatedAt time.Time
	Items     []OrderItem `gorm:"foreignKey:OrderID"`
}

// OrderItem model: copies data from CartItems into Order
type OrderItem struct {
	ID        uint `gorm:"primaryKey;autoIncrement"`
	OrderID   uint `gorm:"index;not null"`
	ProductID uint `gorm:"not null"`
	Quantity  int  `gorm:"not null"`
	Product   Product
}
