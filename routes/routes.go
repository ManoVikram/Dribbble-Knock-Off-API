package routes

import (
	"github.com/ManoVikram/flexibbble-api/handlers"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	// Auth endpoints
	server.POST("/api/auth/signup", handlers.SignupHandler)
	server.POST("/api/auth/login", handlers.LoginHandler)

	// Create project endpoint
	server.POST("/api/createproject", handlers.CreateProjectHandler)
	
	// Select all projects endpoint
	server.GET("/api/allprojects", handlers.FetchAllProjectsHandler)
}