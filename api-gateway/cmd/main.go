package main

import (
	"api-gateway/internal/database"
	"api-gateway/internal/middleware"
	"api-gateway/internal/routes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func init() {
	database.DBConnect()
}

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	r.Use(middleware.CORSMiddleware())

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Public routes
	r.POST("/login", routes.ProxyAuthService)
	r.POST("/signup", routes.ProxyAuthService)
	r.GET("/products/:id", routes.ProxyProductService)
	r.GET("/products", routes.ProxyProductService)
	r.GET("/categories", routes.ProxyProductService)

	// Group for authenticated routes

	auth := r.Group("/")
	auth.Use(middleware.RequireAuth)
	{
		auth.POST("/products", routes.ProxyProductService)
		auth.DELETE("/products/:id", routes.ProxyProductService)

		auth.POST("/orders", routes.ProxyOrderService)
		auth.GET("/orders", routes.ProxyOrderService)
		auth.PUT("/orders/:id", routes.ProxyOrderService)
		auth.DELETE("/orders/:id", routes.ProxyOrderService)

		auth.POST("/logout", routes.ProxyAuthService)

		auth.POST("/categories", routes.ProxyProductService)
		auth.DELETE("/categories", routes.ProxyProductService)

	}

	r.Run(":3000")
}
