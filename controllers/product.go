package controllers

import (
	"fmt"
	"image"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"bogbon-api/models"
	"bogbon-api/repository"
	"github.com/disintegration/imaging"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

// ListProducts godoc
// @Summary List all products
// @Tags Products
// @Produce json
// @Success 200 {array} models.Product
// @Failure 500 {object} map[string]string
// @Router /products [get]
func ListProducts(c *gin.Context) {
	products, err := repository.GetAllProducts()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}

// GetProduct godoc
// @Summary Get a product by ID
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [get]
func GetProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	product, err := repository.GetProductByID(uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, product)
}

// CreateProduct godoc
func CreateProduct(c *gin.Context) {
	var input struct {
		Price       int                    `json:"price"`
		Stock       int                    `json:"stock"`
		Type        string                 `json:"type"`
		Image       string                 `json:"image"`
		Categories  []struct{ ID uint }    `json:"categories"`
		Translations map[string]struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"translations"`
	}

	// Bind the JSON body to the input struct
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Create the product
	product := models.Product{
		Price: input.Price,
		Stock: input.Stock,
		Type:  input.Type,
		Image: input.Image,
	}

	// Set the name and description in the default language (English, "en")
	if input.Translations["en"].Name != "" {
		product.Name = input.Translations["en"].Name
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "English translation for name is required"})
		return
	}
	if input.Translations["en"].Description != "" {
		product.Description = input.Translations["en"].Description
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "English translation for description is required"})
		return
	}

	// Create the product in the database
	if err := repository.CreateProduct(&product); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Now handle the translations for other languages
	for lang, translation := range input.Translations {
		if lang != "en" { // Skip "en" since it was already added as the default translation
			translationRecord := models.ProductTranslation{
				ProductID:    product.ID,
				LanguageCode: lang,
				Name:         translation.Name,
				Description:  translation.Description,
			}
			if err := repository.CreateTranslation(&translationRecord); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save translations"})
				return
			}
		}
	}

	c.JSON(http.StatusCreated, product)
}


// UpdateProduct godoc
// @Summary Update a product
// @Tags Products
// @Produce json
// @Param id path int true "Product ID"
// @Success 200 {object} models.Product
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /products/{id} [put]
func UpdateProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	var input struct {
		Price        int    `json:"price"`
		Stock        int    `json:"stock"`
		Type         string `json:"type"`
		Image        string `json:"image"`
		Translations map[string]struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"translations"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product := models.Product{
		ID:    uint(id),
		Price: input.Price,
		Stock: input.Stock,
		Type:  input.Type,
		Image: input.Image,
	}

	err = repository.UpdateProduct(&product, input.Translations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UploadProductImage handles image upload for a specific product
func UploadProductImage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Image file is required"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, JPEG, or PNG files are allowed"})
		return
	}

	srcFile, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
		return
	}
	defer srcFile.Close()

	srcImage, _, err := image.Decode(srcFile)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
		return
	}

	resizedImage := imaging.Resize(srcImage, 400, 0, imaging.Lanczos)

	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	imageUUID := uuid.New().String()
	filename := fmt.Sprintf("product_%s%s", imageUUID, ext)
	fullPath := filepath.Join(uploadPath, filename)

	// Compress based on format
	switch ext {
	case ".jpg", ".jpeg":
		err = imaging.Save(resizedImage, fullPath, imaging.JPEGQuality(70))
	case ".png":
		err = imaging.Save(resizedImage, fullPath, imaging.PNGCompressionLevel(png.BestCompression))
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
		return
	}

	// Update the product's image in the database
	err = repository.UpdateProductImage(uint(id), fullPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product image"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Image uploaded successfully"})
}

func DeleteProduct(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := repository.DeleteProduct(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
