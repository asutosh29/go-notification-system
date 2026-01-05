package hub

import (
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
)

type SseClient struct {
	Id   string
	Name string

	NotifyChan chan database.Notification
}

type Clients struct {
	data map[string]SseClient
	mu   sync.Mutex
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
