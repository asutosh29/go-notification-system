package api

import (
	"net/http"
	"strconv"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/gin-gonic/gin"
)

func AddNotification(c *gin.Context) {
	var notification database.NotificationResp

	if err := c.Bind(&notification); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input type",
		})
		return
	}

	_, err := database.AddNotification(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't add resource to database",
			"error":   err,
		})
	}

	c.JSON(http.StatusCreated, notification)
}
func AllNotification(c *gin.Context) {
	notifications, err := database.AllNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't fetch resource from database",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
func GetNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input type",
		})
		return
	}

	notification, err := database.GetNotification(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't fetch resource from database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notification)
}
func DeleteNotification(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid input type",
		})
		return
	}

	err = database.DeleteNotification(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't delete resource from database",
			"error":   err,
		})
		return
	}

	c.JSON(http.StatusOK, id)
}
