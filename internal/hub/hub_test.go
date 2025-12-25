package hub

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/asutosh29/go-gin/internal/database"
)

// waitDuration gives the channels time to process messages
const waitDuration = 50 * time.Millisecond

func TestHub_Lifecycle_WithMutex(t *testing.T) {
	h := NewHub()

	go h.Listen()

	u1 := database.User{Id: "user-1", Name: "Alice"}

	t.Log("Adding Client...")
	h.AddClient(u1)

	// Let the hub go routine process the query before checking
	time.Sleep(waitDuration)

	count := h.clients.Count()
	if count != 1 {
		t.Errorf("Expected 1 client, got %d", count)
	} else {
		t.Log("PASS: Client count is 1")
	}

	t.Log("Removing Client...")
	h.RemoveClient(u1)
	time.Sleep(waitDuration)

	count = h.clients.Count()
	if count != 0 {
		t.Errorf("Expected 0 clients, got %d", count)
	} else {
		t.Log("PASS: Client count is 0")
	}
}

func TestHub_Concurrency_Safe(t *testing.T) {
	t.Log("Starting Concurrency Stress Test on Clients Mutex...")

	h := NewHub()
	go h.Listen()

	var wg sync.WaitGroup
	workers := 50

	// Concurrent Adds
	// Launch 50 goroutines. Each adds a user.
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			u := database.User{Id: fmt.Sprintf("u-%d", id)}
			h.AddClient(u)
		}(i)
	}

	wg.Wait()

	// Wait for Hub to process the channel buffer
	time.Sleep(200 * time.Millisecond)

	// Verify Count
	// If Mutex logic in Clients struct is wrong, this might be incorrect.
	finalCount := h.clients.Count()
	if finalCount != workers {
		t.Errorf("Expected %d clients, got %d. (Is the channel blocked?)", workers, finalCount)
	} else {
		t.Logf("PASS: Successfully handled %d concurrent additions", finalCount)
	}
}

func TestHub_Broadcast_Access(t *testing.T) {
	h := NewHub()
	go h.Listen()

	h.AddClient(database.User{Id: "reader", Name: "Reader"})
	time.Sleep(waitDuration)

	// Trigger Broadcast
	// This tests that the range loop in Listen()
	//: "for _, client := range h.clients.Clients().data"
	// doesn't panic.
	t.Log("Sending broadcast...")
	h.Broadcast(database.Notification{NotificationResp: database.NotificationResp{
		Title:       "Test",
		Description: "Hello",
	},
	})

	time.Sleep(waitDuration)
	t.Log("PASS: Broadcast finished without panic")
}
