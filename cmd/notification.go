package main

import (
	"net/http"
	"strconv"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/gin-gonic/gin"
)

func (app *Application) AddNotification(c *gin.Context) {
	var notification database.NotificationResp

	if err := c.Bind(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input type",
		})
		return
	}

	tx := app.db.Create(&database.Notification{
		NotificationResp: notification,
	})
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't add resource to database",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, notification)
}
func (app *Application) AllNotification(c *gin.Context) {
	var notifications []database.Notification
	tx := app.db.Find(&notifications)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't fetch resource from database",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
func (app *Application) GetNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input type",
		})
		return
	}

	var notification database.Notification

	tx := app.db.First(&notification, id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't fetch resource from database",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notification)
}
func (app *Application) DeleteNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input type",
		})
		return
	}

	tx := app.db.Unscoped().Delete(&database.Notification{}, id)
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't delete resource from database",
			"error":   tx.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, id)
}
