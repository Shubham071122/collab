package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/shubham071122/collab/internal/config"
	"github.com/shubham071122/collab/internal/database"
	"github.com/shubham071122/collab/internal/routes"
	"github.com/shubham071122/collab/internal/collaboration"
)

func main() {

	cfg := config.LoadConfig()

	db := database.ConnectPostgres(cfg)

	hub := collaboration.NewHub(db)
	go hub.Run()

	router := gin.Default()

	routes.RegisterRoutes(router, db, hub)

	if err := router.Run(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
