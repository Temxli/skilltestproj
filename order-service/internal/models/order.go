package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	OrderID    string      `gorm:"uniqueIndex;not null" json:"order_id"`
	CustomerID uint        `gorm:"not null" json:"customer_id"`
	Items      []OrderItem `gorm:"foreignKey:OrderID;references:OrderID" json:"items"`
}

type OrderItem struct {
	gorm.Model
	OrderID     string  `gorm:"index;not null" json:"order_id"`
	ProductID   uint    `gorm:"not null" json:"product_id"`
	ProductName string  `gorm:"not null" json:"product_name"`
	UnitPrice   float64 `gorm:"not null" json:"unit_price"`
}

type Product struct {
	gorm.Model
	Name  string  `gorm:"not null" json:"name"`
	Price float64 `gorm:"not null" json:"price"`
}
type User struct {
	gorm.Model
	Name     string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
}
