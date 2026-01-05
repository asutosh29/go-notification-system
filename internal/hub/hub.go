package hub

import (
	"log"
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
)

type Hub struct {
	data             map[string]*SseClient
	connect          chan *SseClient
	disconnect       chan *SseClient
	BroadcastChannel chan database.Notification

	quit chan struct{}
	wg   sync.WaitGroup

	once sync.Once // only for safe closure
}

func NewHub() *Hub {
	h := &Hub{
		data:             make(map[string]*SseClient),
		connect:          make(chan *SseClient, 100),
		disconnect:       make(chan *SseClient, 100),
		BroadcastChannel: make(chan database.Notification, 100),

		quit: make(chan struct{}),
		wg:   sync.WaitGroup{},
	}
	h.wg.Add(1)
	return h

}

func (h *Hub) Listen() {

	defer h.wg.Done()

	log.Println("Hub started listening...")

	for {
		select {

		case user, ok := <-h.connect:
			if !ok {
				return
			}
			h.add(user)
			log.Print("New client connected: ", user.Id)
			log.Print("Num client: ", h.count())
		case user, ok := <-h.disconnect:
			if !ok {
				return
			}
			log.Print("user being disconnected: ", user)
			if client, exists := h.data[user.Id]; exists {
				close(client.NotifyChan)
				h.remove(user)
				log.Print("Client disconnected: ", user.Id)
				log.Print("Num client: ", h.count())
			}
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
		case <-h.quit:
			log.Println("Hub shutting down...")

			for _, client := range h.data {
				close(client.NotifyChan)
			}
			log.Println("Hub shutdown complete. All clients disconnected.")
			return
		}
	}
}

func (h *Hub) close() {
	h.once.Do(func() {
		close(h.quit) // send the signal <-h.quit in the select statement
		h.wg.Wait()
	})
}

func (h *Hub) AddClient(user *SseClient) {
	select {
	case h.connect <- user:
	case <-h.quit:
	}
}

func (h *Hub) RemoveClient(user *SseClient) {
	select {
	case h.disconnect <- user:
	case <-h.quit:
	}
}

func (h *Hub) BroadcastNotification(notif database.Notification) {
	select {
	case h.BroadcastChannel <- notif:
	case <-h.quit:
	}
}

func (h *Hub) add(user *SseClient) {
	h.data[user.Id] = user
}

func (h *Hub) remove(user *SseClient) {
	delete(h.data, user.Id)
}
func (h *Hub) count() int {
	return len(h.data)
}
