package api

import "github.com/gin-gonic/gin"

func InitRouter() *gin.Engine {
	router := gin.Default()
	router.GET("/", getWelcomeMessage)
	router.GET("/health", getHealthStatus)

	notificationGroup := router.Group("/notification")
	{
		notificationGroup.GET("/", AllNotification)
		notificationGroup.POST("/", AddNotification)
		notificationGroup.GET("/:id", GetNotification)
		notificationGroup.DELETE("/:id", DeleteNotification)
	}

	return router
}
