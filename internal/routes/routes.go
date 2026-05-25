package routes

import (
	"github.com/gin-gonic/gin"
	
	"github.com/shubham071122/collab/internal/auth"
)

func RegisterRoutes(router *gin.Engine) {

	authHandler := auth.NewHandler()

	api := router.Group("/api")
	{
		api.GET("/health", healthCheck)

		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.GET("/user", authHandler.User)
			authRoutes.POST("/logout", authHandler.Logout)
		}
	}
}

func healthCheck(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "server running smoothly",
	})
}