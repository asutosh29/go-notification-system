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
