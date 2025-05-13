package repository

import (
	"fmt"
	"sort"
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
	IsOriginal  bool
}

type CategoryFilter struct {
	Q string
}

func FilterProducts(f ProductFilter, includeImages bool) ([]models.Product, error) {
	var products []models.Product

	query := config.DB.Preload("Categories").Preload("Translations").Preload("Images", "is_original = ?", f.IsOriginal)

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

	// If images should be included, sort images after querying
	for i := range products {
		// Sort images array: 'default' images first
		images := products[i].Images
		products[i].Images = sortImagesByDefault(images)
	}

	fmt.Println("THE QUERY: ", query)

	return products, nil
}

// sortImagesByDefault sorts images such that images with 'default' in the URL appear first.
func sortImagesByDefault(images []models.ProductImage) []models.ProductImage {
	sortedImages := make([]models.ProductImage, len(images))
	copy(sortedImages, images)

	sort.SliceStable(sortedImages, func(i, j int) bool {
		isDefaultI := strings.Contains(sortedImages[i].URL, "default")
		isDefaultJ := strings.Contains(sortedImages[j].URL, "default")

		if isDefaultI && !isDefaultJ {
			return true
		}
		if !isDefaultI && isDefaultJ {
			return false
		}

		return sortedImages[i].ID < sortedImages[j].ID
	})

	return sortedImages
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
