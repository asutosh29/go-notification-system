package hub

import (
	"log"
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
)

type hub struct {
	clients    Clients
	mu         sync.Mutex
	connect    chan database.User
	disconnect chan database.User
	broadcast  chan database.Notification
}

func NewHub() *hub {
	return &hub{
		mu: sync.Mutex{},
		clients: Clients{
			data: make(map[string]database.User),
			mu:   sync.Mutex{},
		},
		connect:    make(chan database.User), // TODO: Make them Buffered to avoid blocking due to bad internet speed
		disconnect: make(chan database.User),
		broadcast:  make(chan database.Notification),
	}
}

func (h *hub) Listen() {
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
		case notif := <-h.broadcast:
			log.Print("Broadcasting notification: ", notif)
			for _, client := range h.clients.Clients().data {
				log.Print("Notification sent to ", client.Id)
			}
			// send notifications
		}
	}
}

func (h *hub) AddClient(user database.User) {
	h.connect <- user
}

func (h *hub) RemoveClient(user database.User) {
	h.disconnect <- user
	// handle connection cleanup in the main loop
}

func (h *hub) Broadcast(notif database.Notification) {
	h.broadcast <- notif
}
