package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/wanafiq/feed-api/internal/constants"
	"github.com/wanafiq/feed-api/internal/handlers"
	"github.com/wanafiq/feed-api/internal/middleware"
)

func InitRoutes(m *middleware.Middleware, authHandler *handlers.AuthHandler, userHandler *handlers.UserHandler) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	{
		// Authentication routes
		api.POST("/auth/register", authHandler.Register)
		api.POST("/auth/login", authHandler.Login)
	}

	privateApi := router.Group("/api/v1")
	privateApi.Use(m.RequireAuth())
	{
		// User routes
		privateApi.GET("/users/:userID", userHandler.GetByID)
		privateApi.PUT("/users/:userID/follow", userHandler.Follow)
		privateApi.PUT("/users/:userID/unfollow", userHandler.Unfollow)
		privateApi.PUT("/users/:userID", m.RequireRoles(constants.RoleAdmin), userHandler.Deactivate)
	}

	return router
}
