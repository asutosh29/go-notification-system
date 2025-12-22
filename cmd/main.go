package main

import (
	"fmt"
	"log"

	"github.com/asutosh29/go-gin/database"
	"github.com/asutosh29/go-gin/internal/config"
	"github.com/asutosh29/go-gin/internal/env"
	"github.com/gin-gonic/gin"
)

func main() {
	// Configs
	config := config.InitConfig()

	// DB Setup
	db, err := database.ConnectDb(config)
	if err != nil {
		log.Fatal("Error connecting to Database: ", err)
	}

	// Application setup
	app := Application{
		db:     db,
		config: config,
	}

	// Web server Setup
	PORT := env.GetEnvString("SERVER_PORT", "8080")

	router := gin.Default()
	router.GET("/", app.getWelcomeMessage)
	router.GET("/health", app.getHealthStatus)
	router.Run(fmt.Sprintf("localhost:%s", PORT))
}
