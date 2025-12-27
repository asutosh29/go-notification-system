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
	connect          chan SseClient
	disconnect       chan SseClient
	BroadcastChannel chan database.Notification
}

type HandlerFunc func(*gin.Context)

func NewHub() *Hub {
	return &Hub{
		mu: sync.Mutex{},
		clients: Clients{
			data: make(map[string]SseClient),
			mu:   sync.Mutex{},
		},
		connect:          make(chan SseClient), // TODO: Make them Buffered to avoid blocking due to bad internet speed
		disconnect:       make(chan SseClient),
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
			close(user.NotifyChan)
			h.clients.Remove(user)
			log.Print("Client disconnected: ", user.Id)
			log.Print("Num client: ", h.clients.Count())
			h.mu.Unlock()
		case notif := <-h.BroadcastChannel:
			h.mu.Lock()
			log.Print("Broadcasting notification: ", notif)
			for _, client := range h.clients.Clients().data {
				client.NotifyChan <- notif
				log.Print("Notification sent to ", client.Id)
			}
			h.mu.Unlock()
			// send notifications
		}
	}
}

func (h *Hub) Close() {
	close(h.connect)
	close(h.disconnect)
	close(h.BroadcastChannel)
}

func (h *Hub) AddClient(user SseClient) {
	h.connect <- user
}

func (h *Hub) RemoveClient(user SseClient) {
	h.disconnect <- user
	// handle connection cleanup in the main loop
}

func (h *Hub) BroadcastNotification(notif database.Notification) {
	h.BroadcastChannel <- notif
}
