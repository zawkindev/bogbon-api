package repository

import (
	"bogbon-api/config"
	"bogbon-api/models"
)

func CreateProduct(p *models.Product) error {
	return config.DB.Create(p).Error
}

func GetAllProducts() ([]models.Product, error) {
	var products []models.Product
	err := config.DB.Preload("Categories").Find(&products).Error
	return products, err
}

func GetProductByID(id uint) (*models.Product, error) {
	var p models.Product
	err := config.DB.Preload("Categories").First(&p, id).Error
	return &p, err
}

func UpdateProduct(p *models.Product) error {
	return config.DB.Save(p).Error
}

func DeleteProduct(id uint) error {
	return config.DB.Delete(&models.Product{}, id).Error
}
