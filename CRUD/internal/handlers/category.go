package handlers

import (
	"CRUD/internal/database"
	"CRUD/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateCategory(c *gin.Context) {
	var category models.Category
	if err := c.ShouldBindJSON(&category); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, category)

}

func GetCategories(c *gin.Context) {
	var categories []models.Category
	database.DB.Find(&categories)
	c.JSON(http.StatusOK, categories)

}

func GetCategoryByID(c *gin.Context) {
	catID := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, catID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	c.JSON(http.StatusOK, category)

}

func UpdateCategory(c *gin.Context) {
	catID := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, catID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}

}
func DeleteCategory(c *gin.Context) {
	catID := c.Param("id")
	var category models.Category
	if err := database.DB.First(&category, catID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "category not found"})
		return
	}
	if err := database.DB.Delete(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "category deleted successfully"})
}
