package repository

import (
    "bogbon-api/models"
    "bogbon-api/config"
)

// CreateCategory adds a new category.
func CreateCategory(cat *models.Category) error {
    return config.DB.Create(cat).Error
}

// GetAllCategories returns all categories.
func GetAllCategories() ([]models.Category, error) {
    var cats []models.Category
    err := config.DB.Preload("Products").Find(&cats).Error
    return cats, err
}

// DeleteCategory removes a category by ID.
func DeleteCategory(id uint) error {
    return config.DB.Delete(&models.Category{}, id).Error
}
