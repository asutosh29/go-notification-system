package hub

import (
	"encoding/json"
	"io"
	"log"
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Hub struct {
	clients          Clients
	mu               sync.Mutex
	connect          chan database.User
	disconnect       chan database.User
	BroadcastChannel chan database.Notification
}

type HandlerFunc func(*gin.Context)

func NewHub() *Hub {
	return &Hub{
		mu: sync.Mutex{},
		clients: Clients{
			data: make(map[string]database.User),
			mu:   sync.Mutex{},
		},
		connect:          make(chan database.User), // TODO: Make them Buffered to avoid blocking due to bad internet speed
		disconnect:       make(chan database.User),
		BroadcastChannel: make(chan database.Notification),
	}
}

func (h *Hub) Listen() {
	for {
		select {
		case user := <-h.connect:
			h.mu.Lock()
			h.clients.Add(user)
			log.Print("New client connected: ", user.Id)
			log.Print("Num client: ", h.clients.Count())
			h.mu.Unlock()
		case user := <-h.disconnect:
			h.mu.Lock()
			h.clients.Remove(user)
			log.Print("Client disconnected: ", user.Id)
			log.Print("Num client: ", h.clients.Count())
			h.mu.Unlock()
		case notif := <-h.BroadcastChannel:
			log.Print("Broadcasting notification: ", notif)
			for _, client := range h.clients.Clients().data {
				log.Print("Notification sent to ", client.Id)
			}
			// send notifications
		}
	}
}

func (h *Hub) Close() {
	close(h.connect)
	close(h.disconnect)
	close(h.BroadcastChannel)
}

func (h *Hub) AddClient(user database.User) {
	h.connect <- user
}

func (h *Hub) RemoveClient(user database.User) {
	h.disconnect <- user
	// handle connection cleanup in the main loop
}

func (h *Hub) BroadcastNotification(notif database.Notification) {
	h.BroadcastChannel <- notif
}

// depracated: use controller version
func (h *Hub) StreamHandler() HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Content-Type", "text/event-stream")
		c.Writer.Header().Set("Cache-Control", "no-cache")
		c.Writer.Header().Set("Connection", "keep-alive")
		c.Writer.Header().Set("Transfer-Encoding", "chunked")

		user := database.User{Id: uuid.NewString()}
		h.clients.Add(user)

		defer func() {
			h.clients.Remove(user)
		}()

		c.Stream(func(w io.Writer) bool {
			c.SSEvent("user_connected", user)
			select {
			case notif, ok := <-h.BroadcastChannel:
				if !ok {
					return false // Channel closed
				}
				// Format: "event: <type>\ndata: <json>\n\n"
				jsonNotif, err := json.Marshal(notif)
				if err != nil {
					log.Print("Error unmarshalling notification payload: ", err)
				}

				c.SSEvent("Notification", jsonNotif)
				return true
			case <-c.Request.Context().Done():
				return false // Client disconnected
			}
		})
	}
}
