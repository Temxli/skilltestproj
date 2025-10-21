package database

import (
	"CRUD/internal/models"
)

func Migrate() {
	DB.AutoMigrate(&models.Product{})
	DB.AutoMigrate(&models.Category{})
}
