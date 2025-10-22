package server

import (
	"auth-service/internal/handlers"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"*"}, // Add your frontend URL

		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowHeaders:     []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true, // Enable cookies/auth
	}))

	r.POST("/login", handlers.Login)
	r.POST("/signup", handlers.CreateUser)
	// User routes
	userGroup := r.Group("/users")
	{
		userGroup.POST("/", handlers.CreateUser)
		userGroup.GET("/", handlers.GetUsers)
		userGroup.GET("/:id", handlers.GetUserByID)
		userGroup.PUT("/:id", handlers.UpdateUser)
		userGroup.DELETE("/:id", handlers.DeleteUser)
		userGroup.POST("/logout", handlers.Logout)
		userGroup.PATCH("/:id/password", handlers.ChangePassword)
		userGroup.GET("/me", handlers.GetMyInfo)
	}

	return r
}
