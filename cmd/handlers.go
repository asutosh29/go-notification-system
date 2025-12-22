package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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
