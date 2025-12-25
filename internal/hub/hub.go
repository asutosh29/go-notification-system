package hub

import (
	"log"
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/gin-gonic/gin"
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
