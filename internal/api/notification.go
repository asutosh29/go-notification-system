package api

import (
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/asutosh29/go-gin/internal/hub"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (r *Controller) AddNotification(c *gin.Context) {
	var NotificationResp database.NotificationResp

	if err := c.Bind(&NotificationResp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Could not process input",
			"error":   "Invalid input type",
		})
		return
	}

	notification := database.Notification{
		Title:       NotificationResp.Title,
		Description: NotificationResp.Description,
	}
	_, err := database.AddNotification(notification)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't add resource to database",
			"error":   err.Error(),
		})
		return
	}
	log.Print("Broadcasting notification")
	r.hub.BroadcastNotification(notification)

	c.JSON(http.StatusCreated, notification)
}
func (r *Controller) AllNotification(c *gin.Context) {
	notifications, err := database.AllNotification()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Couldn't fetch resource from database",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, notifications)
}
func (r *Controller) GetNotification(c *gin.Context) {
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
func (r *Controller) DeleteNotification(c *gin.Context) {
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
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, id)
}

func (r *Controller) Stream(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	user := hub.SseClient{
		Id:         uuid.NewString(),
		NotifyChan: make(chan database.Notification, 20),
	}
	r.hub.AddClient(user)
	log.Print("adding id: ", user.Id)

	defer func() {
		r.hub.RemoveClient(user)
	}()

	welcome := false
	c.Stream(func(w io.Writer) bool {
		if !welcome {
			c.SSEvent("user_connected", user.Id)
			welcome = true
			return true
		}
		select {
		case notif, ok := <-user.NotifyChan:
			if !ok {
				return false // Channel closed
			}
			// Format: "event: <type>\ndata: <json>\n\n"
			c.SSEvent("notification", notif)
			return true
		case <-c.Request.Context().Done():
			return false // Client disconnected
		}
	})
}
