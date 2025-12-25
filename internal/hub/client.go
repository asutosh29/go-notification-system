package hub

import (
	"sync"

	"github.com/asutosh29/go-gin/internal/database"
)

type Clients struct {
	data map[string]database.User
	mu   sync.Mutex
}

func (c *Clients) Add(user database.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.data[user.Id] = user
}

func (c *Clients) Remove(user database.User) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.data, user.Id)
}
func (c *Clients) Count() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.data)
}

func (c *Clients) Clients() *Clients {
	return c
}
