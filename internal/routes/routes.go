package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/handlers"
)

func InitRoutes(authHandler *handlers.AuthHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		// Authentication routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	return router
}
