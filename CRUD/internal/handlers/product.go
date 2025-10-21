package handlers

import (
	"CRUD/internal/database"
	"CRUD/internal/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CreateProduct(c *gin.Context) {
	var p models.Product
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}

	if err := database.DB.Create(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create product", "details": err.Error()})
		return
	}

	// reload with category
	if err := database.DB.Preload("Category").First(&p, p.ID).Error; err == nil {
		c.JSON(http.StatusCreated, p)
		return
	}
	c.JSON(http.StatusCreated, p)
}

func GetProducts(c *gin.Context) {
	var p []models.Product
	database.DB.Find(&p)
	c.JSON(http.StatusOK, p)
}
func GetProduct(c *gin.Context) {
	pID := c.Param("id")
	var p models.Product
	if error := database.DB.First(&p, pID).Error; error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}

	c.JSON(http.StatusOK, p)

}
func UpdateProduct(c *gin.Context) {
	pID := c.Param("id")
	var p models.Product
	if error := database.DB.First(&p, pID).Error; error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request", "details": err.Error()})
		return
	}
	if err := database.DB.Save(&p).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update product", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, p)

}
func DeleteProduct(c *gin.Context) {
	pID := c.Param("id")
	var p models.Product
	if error := database.DB.Delete(&p, pID).Error; error != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "product not found"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "product deleted successfully"})
}
