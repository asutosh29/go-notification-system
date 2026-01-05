package hub

import (
	"log"

	"github.com/asutosh29/go-gin/internal/database"
	"github.com/gin-gonic/gin"
)

type Hub struct {
	data             map[string]SseClient
	connect          chan SseClient
	disconnect       chan SseClient
	BroadcastChannel chan database.Notification
}

type HandlerFunc func(*gin.Context)

func NewHub() *Hub {
	return &Hub{
		data:             make(map[string]SseClient),
		connect:          make(chan SseClient, 50), // TODO: Make them Buffered to avoid blocking due to bad internet speed
		disconnect:       make(chan SseClient, 50),
		BroadcastChannel: make(chan database.Notification, 100),
	}
}

func (h *Hub) Listen() {
	for {
		select {
		case user, ok := <-h.connect:
			if !ok {
				return
			}
			h.Add(user)
			log.Print("New client connected: ", user.Id)
			log.Print("Num client: ", h.Count())
		case user, ok := <-h.disconnect:
			if !ok {
				return
			}
			if client, exists := h.data[user.Id]; exists {
				close(client.NotifyChan)
				delete(h.data, user.Id)
				log.Print("Client disconnected: ", user.Id)
			}
			log.Print("Client disconnected: ", user.Id)
			log.Print("Num client: ", h.Count())
		case notif, ok := <-h.BroadcastChannel:
			if !ok {
				return
			}
			log.Print("Broadcasting notification: ", notif)
			for _, client := range h.data {
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
