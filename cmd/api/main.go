package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/shubham071122/collab/internal/config"
	"github.com/shubham071122/collab/internal/database"
	"github.com/shubham071122/collab/internal/routes"
)

func main() {

	cfg := config.LoadConfig()

	db := database.ConnectPostgres(cfg)

	router := gin.Default()

	routes.RegisterRoutes(router, db)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
