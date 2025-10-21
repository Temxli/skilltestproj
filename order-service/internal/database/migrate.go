package database

import (
	"order-service/internal/models"
)

func Migrate() {
	DB.AutoMigrate(&models.Order{})
	DB.AutoMigrate(&models.OrderItem{})
}
