package repository

import (
	"strings"

	"bogbon-api/config"
	"bogbon-api/models"
)

type ProductFilter struct {
	MinPrice    *int
	MaxPrice    *int
	Type        string
	InStock     *bool
	CategoryIDs []uint
	Q           string
}

type CategoryFilter struct {
	Q string
}

// FilterProducts returns products matching any combination of filters
func FilterProducts(f ProductFilter) ([]models.Product, error) {
	db := config.DB.
		Model(&models.Product{}).
		Preload("Categories").
		Preload("Translations")

	// Price range
	if f.MinPrice != nil {
		db = db.Where("price >= ?", *f.MinPrice)
	}
	if f.MaxPrice != nil {
		db = db.Where("price <= ?", *f.MaxPrice)
	}

	// Type
	if f.Type != "" {
		db = db.Where("type = ?", f.Type)
	}

	// In‐stock?
	if f.InStock != nil {
		if *f.InStock {
			db = db.Where("stock > 0")
		} else {
			db = db.Where("stock = 0")
		}
	}

	// By category
	if len(f.CategoryIDs) > 0 {
		db = db.Joins("JOIN category_products cp ON cp.product_id = products.id").
			Where("cp.category_id IN ?", f.CategoryIDs)
	}

	// Full‐text search on translations
	if f.Q != "" {
		like := "%" + strings.ToLower(f.Q) + "%"
		db = db.Joins("JOIN product_translations pt ON pt.product_id = products.id").
			Where("LOWER(pt.name) LIKE ? OR LOWER(pt.description) LIKE ?", like, like).
			Group("products.id")
	}

	var products []models.Product
	err := db.Find(&products).Error
	return products, err
}

// FilterCategories by translation name
func FilterCategories(f CategoryFilter) ([]models.Category, error) {
	db := config.DB.
		Model(&models.Category{}).
		Preload("Translations")

	if f.Q != "" {
		like := "%" + strings.ToLower(f.Q) + "%"
		db = db.Joins("JOIN category_translations ct ON ct.category_id = categories.id").
			Where("LOWER(ct.name) LIKE ?", like).
			Group("categories.id")
	}

	var cats []models.Category
	err := db.Find(&cats).Error
	return cats, err
}
