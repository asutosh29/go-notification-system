package main

import (
	"net/http"

	"github.com/asutosh29/go-gin/internal/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Application struct {
	db     *gorm.DB
	config config.Config
}

// Handlers
func (app *Application) getWelcomeMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Gin server",
	})
}
func (app *Application) getHealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server Healthy",
	})
}
