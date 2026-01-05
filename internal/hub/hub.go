package hub

import (
	"log"
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/gin-gonic/gin"
)

type Hub struct {
	data             map[string]SseClient
	mu               sync.Mutex
	connect          chan SseClient
	disconnect       chan SseClient
	BroadcastChannel chan database.Notification
}

type HandlerFunc func(*gin.Context)

func NewHub() *Hub {
	return &Hub{
		mu:               sync.Mutex{},
		data:             make(map[string]SseClient),
		connect:          make(chan SseClient, 50), // TODO: Make them Buffered to avoid blocking due to bad internet speed
		disconnect:       make(chan SseClient, 50),
		BroadcastChannel: make(chan database.Notification, 100),
	}
}

func (h *Hub) Listen() {
	for {
		select {
		case user := <-h.connect:
			h.Add(user)
			log.Print("New client connected: ", user.Id)
			log.Print("Num client: ", h.Count())
		case user := <-h.disconnect:
			close(user.NotifyChan)
			h.Remove(user)
			log.Print("Client disconnected: ", user.Id)
			log.Print("Num client: ", h.Count())
		case notif := <-h.BroadcastChannel:
			log.Print("Broadcasting notification: ", notif)
			for _, client := range h.Clients().data {
				log.Printf("client %+v", client)
				select {
				case client.NotifyChan <- notif:
					log.Print("Notification sent to ", client.Id)
				default:
					log.Printf("Skipping client %s (buffer full)", client.Id)
				}

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

func (h *Hub) Add(user SseClient) {
	h.data[user.Id] = user
}

func (h *Hub) Remove(user SseClient) {
	delete(h.data, user.Id)
}
func (h *Hub) Count() int {
	return len(h.data)
}

func (h *Hub) Clients() *Hub {
	return h
}
