package models

import (
	"gorm.io/gorm"
	"time"
)

// Category: fixed taxonomy
type Category struct {
	ID       uint      `gorm:"primaryKey;autoIncrement"`
	Name     string    `gorm:"unique;not null"`
	Products []Product `gorm:"many2many:category_products;constraint:OnDelete:CASCADE;"`
}

// Product: can be "plant" or "service" determined by Type
type Product struct {
	ID          uint   `gorm:"primaryKey;autoIncrement"`
	Name        string `gorm:"not null"`
	Description string
	Price       int        `gorm:"not null"`
	Stock       int        `gorm:"not null"`
	Type        string     `gorm:"type:VARCHAR(20);not null"`
	Categories  []Category `gorm:"many2many:category_products;"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// CartItem: belongs to a session
type CartItem struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	SessionID string `gorm:"index;not null"`
	ProductID uint   `gorm:"not null"`
	Quantity  int    `gorm:"not null;default:1"`
}

// Order: created per session
type Order struct {
	ID        uint   `gorm:"primaryKey;autoIncrement"`
	SessionID string `gorm:"index;not null;unique"`
	IsPaid    bool   `gorm:"default:false"`
	CreatedAt time.Time
	Items     []CartItem `gorm:"foreignKey:SessionID;references:SessionID"`
}
