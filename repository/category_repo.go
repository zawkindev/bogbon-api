package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
)

// CreateCategory adds a new category.
func CreateCategory(cat *models.Category) error {
	return config.DB.Create(cat).Error
}

// GetAllCategories returns all categories.
func GetAllCategories() ([]models.Category, error) {
	var cats []models.Category
	err := config.DB.Find(&cats).Error
	return cats, err
}

// DeleteCategory removes a category by ID.
func DeleteCategory(id uint) error {
	return config.DB.Delete(&models.Category{}, id).Error
}

// Update category
func UpdateCategory(id uint, updatedData *models.Category) (*models.Category, error) {
	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		return nil, err
	}

	category.Name = updatedData.Name
	if err := config.DB.Save(&category).Error; err != nil {
		return nil, err
	}

	return &category, nil
}
