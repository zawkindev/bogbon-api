package controllers

import (
	"net/http"
	"strconv"

	"bogbon-api/models"
	"bogbon-api/repository"
	"github.com/gin-gonic/gin"
)

// ListCategories godoc
// @Summary Get all categories
// @Tags Categories
// @Produce json
// @Success 200 {array} models.Category
// @Failure 500 {object} map[string]string
// @Router /api/categories [get]

// ListCategories responds with all categories.
func ListCategories(c *gin.Context) {
	cats, err := repository.GetAllCategories()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, cats)
}

// CreateCategory godoc
// @Summary Create a new category
// @Tags Categories
// @Accept json
// @Produce json
// @Param category body models.Category true "Category to create"
// @Success 201 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories [post]

// CreateCategory adds a new category.
func CreateCategory(c *gin.Context) {
	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := repository.CreateCategory(&input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

// DeleteCategory godoc
// @Summary Delete a category by ID
// @Tags Categories
// @Param id path int true "Category ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [delete]

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

// UpdateCategory godoc
// @Summary Update an existing category
// @Tags Categories
// @Accept json
// @Produce json
// @Param id path int true "Category ID"
// @Param category body models.Category true "Updated category data"
// @Success 200 {object} models.Category
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /categories/{id} [put]

func UpdateCategory(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid category ID"})
		return
	}

	var input models.Category
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	updated, err := repository.UpdateCategory(uint(id), &input)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, updated)
}
