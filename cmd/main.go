package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()
	router.GET("/", getWelcomeMessage)
	router.GET("/health", getHealthStatus)

	router.Run("localhost:8000")
}

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
