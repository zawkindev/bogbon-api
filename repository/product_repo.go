package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
	"errors"
	"os"
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
			ProductID:    p.ID,
			LanguageCode: lang,
			Name:         translation.Name,
			Description:  translation.Description,
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

// UpdateProductImage updates the image URL for a product
func UpdateProductImage(productID uint, imagePath string) error {
	var product models.Product
	if err := config.DB.First(&product, productID).Error; err != nil {
		return err
	}

	baseURL := os.Getenv("BASE_URL")
	product.Image = baseURL + "/" + imagePath
	return config.DB.Save(&product).Error
}

func CreateTranslation(translation *models.ProductTranslation) error {
	return config.DB.Create(translation).Error
}

func GetAllProductsWithTranslations() ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Preload("Categories").
		Preload("Translations").
		Find(&products).Error
	if err != nil {
		return nil, err
	}
	return products, nil
}

func GetProductWithTranslations(id uint) (*models.Product, error) {
	var product models.Product
	err := config.DB.
		Preload("Categories").
		Preload("Translations").
		First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func UpdateProduct(product *models.Product, translations map[string]struct {
	Name        string
	Description string
}) error {
	tx := config.DB.Begin()

	// Update product fields
	if err := tx.Model(&models.Product{}).Where("id = ?", product.ID).
		Updates(models.Product{
			Price: product.Price,
			Stock: product.Stock,
			Type:  product.Type,
			Image: product.Image,
		}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Remove existing translations
	if err := tx.Where("product_id = ?", product.ID).Delete(&models.ProductTranslation{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Add new translations
	for lang, t := range translations {
		translation := models.ProductTranslation{
			ProductID:    product.ID,
			LanguageCode: lang,
			Name:         t.Name,
			Description:  t.Description,
		}
		if err := tx.Create(&translation).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	// Commit the transaction
	return tx.Commit().Error
}

func DeleteProduct(id uint) error {
	tx := config.DB.Begin()

	// Delete associated translations first
	if err := tx.Where("product_id = ?", id).Delete(&models.ProductTranslation{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Delete the product
	if err := tx.Delete(&models.Product{}, id).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit the transaction
	return tx.Commit().Error
}
