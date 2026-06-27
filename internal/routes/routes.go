package routes

import (
	"database/sql"

	"github.com/gin-gonic/gin"

	"github.com/shubham071122/collab/internal/middleware"
	"github.com/shubham071122/collab/internal/response"

	"github.com/shubham071122/collab/internal/auth"
	"github.com/shubham071122/collab/internal/collaboration"
	"github.com/shubham071122/collab/internal/config"
	"github.com/shubham071122/collab/internal/project"
	"github.com/shubham071122/collab/internal/subscription"
	"github.com/shubham071122/collab/internal/user"
)

func RegisterRoutes(router *gin.Engine, db *sql.DB, hub *collaboration.Hub, cfg *config.Config) {

	authRepo := auth.NewRepository(db)
	userRepo := user.NewRepository(db)
	projectRepo := project.NewRepository(db)
	subscriptionRepo := subscription.NewRepository(db)

	authService := auth.NewService(authRepo, subscriptionRepo)
	userService := user.NewService(userRepo)
	subscriptionService := subscription.NewService(subscriptionRepo, userRepo, cfg)
	projectService := project.NewService(projectRepo, userRepo, hub, subscriptionService)

	authHandler := auth.NewHandler(authService)
	userHandler := user.NewHandler(userService)
	projectHandler := project.NewHandler(projectService)
	subscriptionHandler := subscription.NewHandler(subscriptionService)

	collaborationHandler := collaboration.NewHandler(hub)

	router.GET("/health", healthCheck)
	api := router.Group("/api/v1")
	{
		// PLANS ROUTE
		api.GET("/plans", subscriptionHandler.GetPlans)

		// AUTH ROUTES
		authRoutes := api.Group("/auth")
		{
			authRoutes.POST("/register", authHandler.Register)
			authRoutes.POST("/login", authHandler.Login)
			authRoutes.POST("/logout", authHandler.Logout)
			authRoutes.POST("/verify-otp", authHandler.VerifyOTP)
			authRoutes.POST("/resend-otp", authHandler.ResendOTP)
		}

		// SUBSCRIPTION ROUTES
		subscriptionRoutes := api.Group("/subscription")
		subscriptionRoutes.Use(middleware.AuthMiddleware())
		{
			subscriptionRoutes.GET("/", subscriptionHandler.GetSubscription)
			subscriptionRoutes.POST("/checkout", subscriptionHandler.CreateCheckout)
			subscriptionRoutes.POST("/verify", subscriptionHandler.VerifyPayment)
			subscriptionRoutes.POST("/cancel", subscriptionHandler.CancelCheckout)
			subscriptionRoutes.GET("/transactions", subscriptionHandler.GetTransactions)
		}
		api.POST("/subscription/webhook", subscriptionHandler.Webhook)

		// USER ROUTES
		userRoutes := api.Group("/user")
		userRoutes.Use(middleware.AuthMiddleware())
		{
			userRoutes.GET("/:id", userHandler.CurrentUser)
		}

		// PROJECT ROUTES
		projectRoutes := api.Group("/project")
		projectRoutes.Use(middleware.AuthMiddleware())
		{
			projectRoutes.GET("/", projectHandler.GetProjects)
			projectRoutes.GET("/:id", projectHandler.GetProject)
			projectRoutes.POST("/", projectHandler.CreateProject)
			projectRoutes.PATCH("/:id", projectHandler.UpdateProject)
			projectRoutes.DELETE("/:id", projectHandler.DeleteProject)

			projectRoutes.POST("/:id/share", projectHandler.ShareProject)
			projectRoutes.GET("/:id/collaborators", projectHandler.GetCollaborators)
			projectRoutes.GET("/:id/members", projectHandler.GetProjectMembers)
			projectRoutes.PATCH("/:id/collaborators/:userId", projectHandler.UpdateCollaboratorPermission)
			projectRoutes.DELETE("/:id/collaborators/:userId", projectHandler.RemoveCollaborator)
		}

		collaborationRoutes := api.Group("/collaboration")
		collaborationRoutes.Use(middleware.AuthMiddleware())
		{
			collaborationRoutes.GET(
				"/project/:id",
				collaborationHandler.Connect,
			)
		}
	}
}

func healthCheck(c *gin.Context) {
	response.JSON(c, response.StatusOK, "API is healthy", nil, nil)
}
