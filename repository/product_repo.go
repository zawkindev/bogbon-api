package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
	"errors"
)

// CreateProduct creates a product and its translations
func CreateProduct(p *models.Product, translations map[string]struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}) (*models.Product, error) {
	// Create the product
	if err := config.DB.Create(p).Error; err != nil {
		return nil, err
	}

	// Now create translations for the product
	for lang, translation := range translations {
		translationRecord := models.ProductTranslation{
			ProductID:   p.ID,
			LanguageCode: lang,
			Name:        translation.Name,
			Description: translation.Description,
		}
		if err := config.DB.Create(&translationRecord).Error; err != nil {
			return nil, err
		}
	}

	// Reload product with translations
	if err := config.DB.Preload("Categories").Preload("Translations").First(p, p.ID).Error; err != nil {
		return nil, err
	}

	return p, nil
}

// GetAllProducts fetches all products with their translations and categories
func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Preload("Categories").Preload("Translations").Find(&products).Error
	return products, err
}

// GetProductByID fetches a product by its ID with its translations and categories
func GetProductByID(id uint) (*models.Product, error) {
	var p models.Product
	err := config.DB.Preload("Categories").Preload("Translations").First(&p, id).Error
	if err != nil {
		if err.Error() == "record not found" {
			return nil, errors.New("product not found")
		}
		return nil, err
	}
	return &p, nil
}

// UpdateProduct updates a product and its translations
func UpdateProduct(p *models.Product, translations map[string]struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}) error {
	// Update the product
	if err := config.DB.Save(p).Error; err != nil {
		return err
	}

	// Delete existing translations
	if err := config.DB.Where("product_id = ?", p.ID).Delete(&models.ProductTranslation{}).Error; err != nil {
		return err
	}

	// Add new translations
	for lang, translation := range translations {
		translationRecord := models.ProductTranslation{
			ProductID:   p.ID,
			LanguageCode: lang,
			Name:        translation.Name,
			Description: translation.Description,
		}
		if err := config.DB.Create(&translationRecord).Error; err != nil {
			return err
		}
	}

	return nil
}

// UpdateProductImage updates the image URL for a product
func UpdateProductImage(productID uint, imagePath string) error {
	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return err
	}

	product.Image = imagePath
	return config.DB.Save(&product).Error
}

// DeleteProduct deletes a product by its ID
func DeleteProduct(id uint) error {
	return config.DB.Delete(&models.Product{}, id).Error
}

func CreateTranslation(translation *models.ProductTranslation) error {
	return config.DB.Create(translation).Error
}

