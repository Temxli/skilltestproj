package server

import (
	"auth-service/internal/handlers"
	"auth-service/internal/middleware"
	"net/http"

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

	r.POST("/login", handlers.Login)
	r.POST("/register", handlers.CreateUser)
	// User routes
	userGroup := r.Group("/users")
	{
		userGroup.GET("/", middleware.RequireAuth, handlers.GetUsers)
		userGroup.GET("/:id", middleware.RequireAuth, handlers.GetUserByID)
		userGroup.PUT("/:id", middleware.RequireAuth, handlers.UpdateUser)
		userGroup.DELETE("/:id", middleware.RequireAuth, handlers.DeleteUser)
		userGroup.POST("/logout", middleware.RequireAuth, handlers.Logout)
		userGroup.PATCH("/:id/password", middleware.RequireAuth, handlers.ChangePassword)
		userGroup.GET("/me", middleware.RequireAuth, handlers.GetMyInfo)
	}

	return r
}
