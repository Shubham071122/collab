package main

import (
	"github.com/gin-gonic/gin"

	"github.com/shubham071122/collab/internal/config"
	"github.com/shubham071122/collab/internal/routes"
)

func main() {

	cfg := config.LoadConfig()

	router := gin.Default()

	routes.RegisterRoutes(router)

	router.Run(": " + cfg.Port)
}