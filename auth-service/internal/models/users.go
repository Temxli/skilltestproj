package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}

type LoginRequest struct {
	Name     string `json:"name" binding:"required"`
	Password string `json:"password" binding:"required"`
}
