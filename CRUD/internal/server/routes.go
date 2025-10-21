package server

import (
	"net/http"

	"CRUD/internal/handlers"
	"CRUD/internal/middleware"

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

	r.GET("/products/:id", handlers.GetProduct)
	r.GET("/products", middleware.RequireAuth, handlers.GetProducts)
	r.POST("/products", handlers.CreateProduct)
	r.PUT("/products/:id", handlers.UpdateProduct)
	r.DELETE("/products/:id", handlers.DeleteProduct)

	// category routes
	r.POST("/categories", handlers.CreateCategory)
	r.GET("/categories", handlers.GetCategories)
	r.GET("/categories/:id", handlers.GetCategoryByID)
	r.PUT("/categories/:id", handlers.UpdateCategory)
	r.DELETE("/categories/:id", handlers.DeleteCategory)
	return r
}
