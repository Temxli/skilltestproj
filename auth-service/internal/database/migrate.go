package database

import (
	"auth-service/internal/models"
)

func Migrate() {
	DB.AutoMigrate(&models.User{})
}
