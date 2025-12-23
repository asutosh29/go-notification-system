package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handlers
func getWelcomeMessage(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Welcome to Gin server",
	})
}
func getHealthStatus(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "Server Healthy",
	})
}
