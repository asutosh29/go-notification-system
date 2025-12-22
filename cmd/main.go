package main

import (
	"fmt"
	"log"

	"github.com/asutosh29/go-gin/internal/config"
	"github.com/asutosh29/go-gin/internal/env"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Configs
	config := config.InitConfig()

	// DB Setup
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", config.DbHost, config.DbUser, config.DbPassword, config.DbName, config.DbPort)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Error connecting to Database: ", err)
	}

	// Application setup
	app := Application{
		db:     db,
		config: *config,
	}

	// Web server Setup
	PORT := env.GetEnvString("SERVER_PORT", "8080")

	router := gin.Default()
	router.GET("/", app.getWelcomeMessage)
	router.GET("/health", app.getHealthStatus)
	router.Run(fmt.Sprintf("localhost:%s", PORT))
}
