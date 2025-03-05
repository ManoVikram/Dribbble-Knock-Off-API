package routes

import (
	"github.com/ManoVikram/flexibbble-api/handlers"
	"github.com/ManoVikram/flexibbble-api/middlewares"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Auth endpoints
	server.POST("/api/auth/signup", handlers.SignupHandler)
	server.POST("/api/auth/login", handlers.LoginHandler)

	protectedRoutes := server.Group("/api")
	protectedRoutes.Use(middlewares.AuthMiddleware())

	// Create project endpoint
	protectedRoutes.POST("/createproject", handlers.CreateProjectHandler)
	
	// Select all projects endpoint
	protectedRoutes.GET("/allprojects", handlers.FetchAllProjectsHandler)
}