package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
)

func CreateCategory(c *models.Category, translations map[string]struct {
	Name string `json:"name"`
}) (*models.Category, error) {
	// Create category
	if err := config.DB.Create(c).Error; err != nil {
		return nil, err
	}

	// Create translations
	for lang, trans := range translations {
		record := models.CategoryTranslation{
			CategoryID:   c.ID,
			LanguageCode: lang,
			Name:         trans.Name,
		}
		if err := config.DB.Create(&record).Error; err != nil {
			return nil, err
		}
	}

	// Reload with translations
	if err := config.DB.Preload("Translations").First(c, c.ID).Error; err != nil {
		return nil, err
	}

	return c, nil
}

// GetAllCategories returns all categories.
func GetAllCategories() ([]models.Category, error) {
	var cats []models.Category
	err := config.DB.Preload("Translations").Find(&cats).Error
	return cats, err
}

// DeleteCategory removes a category by ID.
func DeleteCategory(id uint) error {
	return config.DB.Delete(&models.Category{}, id).Error
}

// Update category
func UpdateCategory(c *models.Category, translations map[string]struct {
	Name string `json:"name"`
}) error {
	// Save the base category (not much to update unless you add more fields)
	if err := config.DB.Save(c).Error; err != nil {
		return err
	}

	// Remove existing translations
	if err := config.DB.Where("category_id = ?", c.ID).Delete(&models.CategoryTranslation{}).Error; err != nil {
		return err
	}

	// Add new translations
	for lang, trans := range translations {
		record := models.CategoryTranslation{
			CategoryID:   c.ID,
			LanguageCode: lang,
			Name:         trans.Name,
		}
		if err := config.DB.Create(&record).Error; err != nil {
			return err
		}
	}

	// Refresh with translations
	return config.DB.Preload("Translations").First(c, c.ID).Error
}
