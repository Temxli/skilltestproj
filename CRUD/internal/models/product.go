package models

import (
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name       string    `gorm:"not null" json:"name"`
	Price      float64   `json:"price"`
	CategoryID *uint     `json:"category_id,omitempty"`
	Category   *Category `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category,omitempty"`
}

type Category struct {
	gorm.Model
	Name     string    `gorm:"not null" json:"name"`
	Products []Product `gorm:"foreignKey:CategoryID" json:"products,omitempty"`
}
type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}
