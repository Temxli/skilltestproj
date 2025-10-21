package server

import (
	"net/http"
	"order-service/internal/handlers"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"}, // Add your frontend URL
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.GET("/orders", handlers.GetAllOrdersHandler)
	r.GET("/orders/:id", handlers.GetOrderHandler)
	r.POST("/orders", handlers.CreateOrderHandler)
	r.PUT("/orders/:id", handlers.UpdateOrderHandler)
	r.DELETE("/orders/:id", handlers.DeleteOrderHandler)

	return r
}
