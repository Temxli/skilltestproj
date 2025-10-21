package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"order-service/internal/database"
	"order-service/internal/models"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// fetchProduct делает HTTP GET к product-service и возвращает имя/цену
func fetchProduct(productID uint) (*struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}, error) {
	base := os.Getenv("PRODUCT_SERVICE_URL")
	if base == "" {
		base = "http://localhost:8080"
	}
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()

	url := fmt.Sprintf("%s/products/%d", base, productID)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("product service returned %d", resp.StatusCode)
	}

	var p struct {
		ID    uint    `json:"id"`
		Name  string  `json:"name"`
		Price float64 `json:"price"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return nil, err
	}
	return &p, nil
}

func CreateOrderHandler(c *gin.Context) {
	var o models.Order

	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if o.OrderID == "" {
		o.OrderID = fmt.Sprintf("ord-%d", time.Now().UnixNano())
	}

	// validate items and fill unit prices from product-service
	for i := range o.Items {
		if o.Items[i].ProductID == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("items[%d].product_id is required", i)})
			return
		}

		p, err := fetchProduct(o.Items[i].ProductID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("product %d validation failed: %v", o.Items[i].ProductID, err)})
			return
		}

		o.Items[i].ProductName = p.Name
		o.Items[i].UnitPrice = p.Price
		o.Items[i].OrderID = o.OrderID
	}

	if err := database.DB.Transaction(func(tx *gorm.DB) error {
		// не создавать ассоциации автоматически	
		if err := tx.Omit("Items").Create(&o).Error; err != nil {
			return err
		}
		if len(o.Items) > 0 {
			// на всякий случай обнулим ID, чтобы БД проставила auto-increment
			for i := range o.Items {
				o.Items[i].ID = 0
			}
			if err := tx.Create(&o.Items).Error; err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, o)
}
func GetOrderHandler(c *gin.Context) {
	orderID := c.Param("id")
	var o models.Order
	if err := database.DB.First(&o, "order_id = ?", orderID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Order not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}

func GetAllOrdersHandler(c *gin.Context) {
	var orders []models.Order
	if err := database.DB.Find(&orders).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, orders)

}

func DeleteOrderHandler(c *gin.Context) {
	orderID := c.Param("id")
	if err := database.DB.Delete(&models.Order{}, "order_id = ?", orderID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})

}
func UpdateOrderHandler(c *gin.Context) {
	orderID := c.Param("id")
	var o models.Order
	if err := c.ShouldBindJSON(&o); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := database.DB.Model(&models.Order{}).Where("order_id = ?", orderID).Updates(o).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, o)
}
