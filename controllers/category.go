package controllers

import (
	"bogbon-api/models"
	"bogbon-api/repository"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// ListCategories responds with all categories (including translations),
// optionally filtered by translation name using ?q=
func ListCategories(c *gin.Context) {
	q := c.Query("q")
	var cats []models.Category
	var err error

	if q != "" {
		// Filter by translated name
		cats, err = repository.FilterCategories(repository.CategoryFilter{Q: q})
	} else {
		// Return all categories
		cats, err = repository.GetAllCategories()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, cats)
}

// CreateCategory adds a new category with translations.
func CreateCategory(c *gin.Context) {
	var input struct {
		Translations map[string]struct {
			Name string `json:"name"`
		} `json:"translations"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	en, ok := input.Translations["en"]
	if !ok || en.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "English translation is required"})
		return
	}

	category := models.Category{}

	createdCategory, err := repository.CreateCategory(&category, input.Translations)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdCategory)
}

// DeleteCategory removes a category by ID.
func DeleteCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid ID"})
		return
	}
	if err := repository.DeleteCategory(uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.Status(http.StatusNoContent)
}

// UpdateCategory updates a category and its translations.
func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var input struct {
		Translations map[string]struct {
			Name string `json:"name"`
		} `json:"translations"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Make sure English name is present
	en, ok := input.Translations["en"]
	if !ok || en.Name == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "English translation is required"})
		return
	}

	updatedCategory := models.Category{ID: uint(id)}

	if err := repository.UpdateCategory(&updatedCategory, input.Translations); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updatedCategory)
}
