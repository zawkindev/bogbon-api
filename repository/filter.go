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
	IsOriginal  *bool
}

type CategoryFilter struct {
	Q string
}

// FilterProducts fetches products with optional filters, and includes images if requested
func FilterProducts(f ProductFilter, includeImages bool) ([]models.Product, error) {
	var products []models.Product
	query := config.DB.Preload("Categories").Preload("Translations")

	// Apply filters to the query
	if f.MinPrice != nil {
		query = query.Where("price >= ?", f.MinPrice)
	}
	if f.MaxPrice != nil {
		query = query.Where("price <= ?", f.MaxPrice)
	}
	if f.Type != "" {
		query = query.Where("type = ?", f.Type)
	}
	if f.InStock != nil {
		query = query.Where("stock > ?", 0)
	}
	if f.IsOriginal != nil {
		query = query.Where("isOriginal = ?", *f.IsOriginal)
	}
	if len(f.CategoryIDs) > 0 {
		query = query.Joins("JOIN category_products ON category_products.product_id = products.id").
			Where("category_products.category_id IN ?", f.CategoryIDs)
	}
	if f.Q != "" {
		query = query.Where("name LIKE ?", "%"+f.Q+"%")
	}

	// Execute the query
	if err := query.Find(&products).Error; err != nil {
		return nil, err
	}

	// If images should be included, preload images
	if includeImages {
		for i := range products {
			if err := config.DB.Preload("Images").Find(&products[i]).Error; err != nil {
				return nil, err
			}
		}
	}

	return products, nil
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
