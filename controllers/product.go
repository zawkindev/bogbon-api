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

// ListProducts now supports ?min_price=&max_price=&type=&in_stock=&category=&q=&include_images=true
func ListProducts(c *gin.Context) {
	var f repository.ProductFilter

	// min_price
	if v := c.Query("min_price"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			f.MinPrice = &i
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid min_price"})
			return
		}
	}
	// max_price
	if v := c.Query("max_price"); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			f.MaxPrice = &i
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid max_price"})
			return
		}
	}
	// type
	f.Type = c.Query("type")

	// in_stock
	if v := c.Query("in_stock"); v != "" {
		b, err := strconv.ParseBool(v)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid in_stock"})
			return
		}
		f.InStock = &b
	}

	// category (can be repeated)
	for _, v := range c.QueryArray("category") {
		if v == "" {
			continue
		}
		if id, err := strconv.Atoi(v); err == nil {
			f.CategoryIDs = append(f.CategoryIDs, uint(id))
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid category id"})
			return
		}
	}

	// search term
	f.Q = c.Query("q")

	// fetch products, including images
	includeImages := c.DefaultQuery("include_images", "false") != "true"
	products, err := repository.FilterProducts(f, includeImages)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, products)
}


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
		Price        int                 `json:"price"`
		Stock        int                 `json:"stock"`
		Type         string              `json:"type"`
		Categories   []struct{ ID uint } `json:"categories"`
		Translations map[string]struct {
			Name        string `json:"name"`
			Description string `json:"description"`
			ShortInfo   string `json:"short_info"`
		} `json:"translations"`
	}

	// Bind JSON
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate English translation
	en, ok := input.Translations["en"]
	if !ok || en.Name == "" || en.Description == "" || en.ShortInfo == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "English name, short_info and description are required"})
		return
	}

	// Create base product model
	product := models.Product{
		Price: input.Price,
		Stock: input.Stock,
		Type:  input.Type,
	}

	// Attach categories
	for _, c := range input.Categories {
		product.Categories = append(product.Categories, models.Category{ID: c.ID})
	}

	// Use the repository to create the product with translations
	createdProduct, err := repository.CreateProduct(&product, input.Translations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdProduct)
}

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
		Translations map[string]struct {
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"translations"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Convert input.Translations to a map of the format expected by the repository
	translations := make(map[string]struct {
		Name        string
		Description string
	})

	for lang, t := range input.Translations {
		translations[lang] = struct {
			Name        string
			Description string
		}{Name: t.Name, Description: t.Description}
	}

	// Create product object
	product := models.Product{
		ID:    uint(id),
		Price: input.Price,
		Stock: input.Stock,
		Type:  input.Type,
	}

	// Pass to repository to update
	err = repository.UpdateProduct(&product, translations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

// UploadProductImage handles image upload for a specific product (multiple images)
func UploadProductImage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Parse the multipart form (to handle file uploads)
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Get the uploaded files
	files := c.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one image file is required"})
		return
	}

	// Create upload directory if it doesn't exist
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Loop over each uploaded file
	var imagePaths []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, JPEG, or PNG files are allowed"})
			return
		}

		// Open the file
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer srcFile.Close()

		// Decode image
		srcImage, _, err := image.Decode(srcFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			return
		}

		// Resize the image
		resizedImage := imaging.Resize(srcImage, 0, 273, imaging.Lanczos)

		// Generate a unique filename
		imageUUID := uuid.New().String()
		filename := fmt.Sprintf("product_%s%s", imageUUID, ext)
		fullPath := filepath.Join(uploadPath, filename)

		// Compress based on file type
		switch ext {
		case ".jpg", ".jpeg":
			err = imaging.Save(resizedImage, fullPath, imaging.JPEGQuality(70))
		case ".png":
			err = imaging.Save(resizedImage, fullPath, imaging.PNGCompressionLevel(png.BestCompression))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
			return
		}

		// Add the image path to the list
		imagePaths = append(imagePaths, fullPath)
	}

	// Now associate the images with the product in the database
	for _, path := range imagePaths {
		err = repository.UpdateProductImage(uint(id), path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product image"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Images uploaded successfully"})
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

func DeleteImage(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid product ID"})
		return
	}

	if err := repository.DeleteProductImage(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

func UpdateProductImageByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
		return
	}

	// Parse the multipart form (to handle file uploads)
	err = c.Request.ParseMultipartForm(10 << 20) // 10 MB limit
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse multipart form"})
		return
	}

	// Get the uploaded files
	files := c.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "At least one image file is required"})
		return
	}

	// Create upload directory if it doesn't exist
	uploadPath := "./uploads"
	if err := os.MkdirAll(uploadPath, os.ModePerm); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to create upload directory"})
		return
	}

	// Loop over each uploaded file
	var imagePaths []string
	for _, file := range files {
		ext := strings.ToLower(filepath.Ext(file.Filename))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Only JPG, JPEG, or PNG files are allowed"})
			return
		}

		// Open the file
		srcFile, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open uploaded file"})
			return
		}
		defer srcFile.Close()

		// Decode image
		srcImage, _, err := image.Decode(srcFile)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid image format"})
			return
		}

		// Resize the image
		resizedImage := imaging.Resize(srcImage, 400, 0, imaging.Lanczos)

		// Generate a unique filename
		imageUUID := uuid.New().String()
		filename := fmt.Sprintf("product_%s%s", imageUUID, ext)
		fullPath := filepath.Join(uploadPath, filename)

		// Compress based on file type
		switch ext {
		case ".jpg", ".jpeg":
			err = imaging.Save(resizedImage, fullPath, imaging.JPEGQuality(70))
		case ".png":
			err = imaging.Save(resizedImage, fullPath, imaging.PNGCompressionLevel(png.BestCompression))
		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported file type"})
			return
		}

		// Add the image path to the list
		imagePaths = append(imagePaths, fullPath)
	}

	// Now associate the images with the product in the database
	for _, path := range imagePaths {
		err = repository.UpdateProductImage(uint(id), path)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product image"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Images uploaded successfully"})
}
