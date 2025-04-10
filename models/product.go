package models

import (
	"time"

	"gorm.io/gorm"
)

type Product struct {
	ID         string         `json:"id" gorm:"primarykey"`
	Name       string         `json:"name"`
	Desciption string         `json:"description"`
	Price      int            `json:"price"`
	Stock      int            `json:"stock"`
	CategoryID uint           `json:"category_id"`
	TagID      uint           `json:"tag_id"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}
