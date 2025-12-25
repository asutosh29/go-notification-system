package api

import (
	"github.com/asutosh29/go-gin/internal/hub"
	"github.com/gin-gonic/gin"
)

type Controller struct {
	hub *hub.Hub
}

func NewController(h *hub.Hub) *Controller {
	return &Controller{
		hub: h,
	}
}

func InitRouter(h *hub.Hub) *gin.Engine {
	router := gin.Default()
	controller := NewController(h)
	router.LoadHTMLGlob("templates/*")
	router.GET("/test", controller.RenderIndex)
	router.GET("/", getWelcomeMessage)
	router.GET("/health", getHealthStatus)

	notificationGroup := router.Group("/notification")
	{
		notificationGroup.GET("/", controller.AllNotification)
		// notificationGroup.GET("/stream", gin.HandlerFunc(h.StreamHandler()))
		notificationGroup.GET("/stream", controller.Stream)

		notificationGroup.POST("/", controller.AddNotification)
		notificationGroup.GET("/:id", controller.GetNotification)
		notificationGroup.DELETE("/:id", controller.DeleteNotification)
	}

	return router
}
