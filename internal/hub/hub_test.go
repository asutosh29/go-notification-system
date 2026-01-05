package hub

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/asutosh29/go-gin/internal/database"
)

// Helper to create a dummy notification
func mockNotification(title string) database.Notification {
	return database.Notification{
		Title:       title,
		Description: "Test Description",
	}
}

// Test 1: Basic Connection and Broadcast
// Ensures normal flow works: Join -> Listen -> Receive -> Leave
func TestHub_BasicBroadcast(t *testing.T) {
	h := NewHub()
	go h.Listen()
	defer h.close()

	// 1. Create a client
	clientID := "client-1"
	notifyChan := make(chan database.Notification, 5)
	client := SseClient{
		Id:         clientID,
		NotifyChan: notifyChan,
	}

	// 2. Add Client
	h.AddClient(&client)

	// Allow a tiny moment for the Hub goroutine to process the add
	time.Sleep(10 * time.Millisecond)

	// 3. Broadcast
	expectedTitle := "Hello World"
	h.BroadcastNotification(mockNotification(expectedTitle))

	// 4. Verify Receipt
	select {
	case notif := <-notifyChan:
		if notif.Title != expectedTitle {
			t.Errorf("Expected title %s, got %s", expectedTitle, notif.Title)
		}
	case <-time.After(1 * time.Second):
		t.Error("Timed out waiting for notification")
	}

	// 5. Remove Client
	h.RemoveClient(&client)
}

// Test 2: Slow Client / Full Buffer
// Ensures that if one client is stuck, the Hub stays alive for others.
func TestHub_SlowClient_DoesNotBlock(t *testing.T) {
	h := NewHub()
	go h.Listen()
	defer h.close()

	// Client A: Small buffer, we will fill it up
	clientA := SseClient{
		Id:         "slow-client",
		NotifyChan: make(chan database.Notification, 1), // Only holds 1
	}

	// Client B: Normal client
	clientB := SseClient{
		Id:         "fast-client",
		NotifyChan: make(chan database.Notification, 5),
	}

	h.AddClient(&clientA)
	h.AddClient(&clientB)
	time.Sleep(10 * time.Millisecond)

	// Fill Client A's buffer manually so it blocks
	clientA.NotifyChan <- mockNotification("filler")

	// Now Broadcast a new message
	// If the Hub logic is wrong, this will BLOCK here waiting for Client A
	h.BroadcastNotification(mockNotification("real-message"))

	// Verify Client B still got it immediately
	select {
	case msg := <-clientB.NotifyChan:
		if msg.Title != "real-message" {
			t.Errorf("Client B got wrong message")
		}
	case <-time.After(100 * time.Millisecond):
		t.Error("Hub blocked! Client B did not receive message because Client A was full.")
	}
}

// Test 3: Double Disconnect Safety
// Ensures calling RemoveClient twice doesn't panic
func TestHub_DoubleDisconnect_NoPanic(t *testing.T) {
	h := NewHub()
	go h.Listen()
	defer h.close()

	client := SseClient{
		Id:         "panic-test-user",
		NotifyChan: make(chan database.Notification, 1),
	}

	h.AddClient(&client)
	time.Sleep(10 * time.Millisecond)

	// First disconnect (Normal)
	h.RemoveClient(&client)

	// Wait for processing
	time.Sleep(10 * time.Millisecond)

	// Second disconnect (Should be ignored safely)
	// If logic is bad, this will panic: "close of closed channel"
	func() {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("Code panicked on double disconnect: %v", r)
			}
		}()
		h.RemoveClient(&client)
	}()

	// Ensure Hub is still alive by trying to broadcast
	h.BroadcastNotification(mockNotification("check-alive"))
}

// Test 4: Heavy Concurrency (Race Detector)
// Simulates many users joining/leaving/broadcasting simultaneously
func TestHub_ConcurrencyLoad(t *testing.T) {
	h := NewHub()
	go h.Listen()
	defer h.close()

	var wg sync.WaitGroup
	userCount := 100

	// Spawn 100 users connecting and disconnecting
	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			c := SseClient{
				Id:         fmt.Sprintf("user-%d", id),
				NotifyChan: make(chan database.Notification, 10),
			}

			h.AddClient(&c)

			// Random small sleep to simulate connection time
			time.Sleep(time.Millisecond * 5)

			h.RemoveClient(&c)
		}(i)
	}

	// Spawn a broadcaster running in parallel
	wg.Add(1)
	go func() {
		defer wg.Done()
		for i := 0; i < 50; i++ {
			h.BroadcastNotification(mockNotification(fmt.Sprintf("msg-%d", i)))
			time.Sleep(time.Millisecond * 2)
		}
	}()

	// Wait for everyone to finish
	done := make(chan bool)
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// Success
	case <-time.After(3 * time.Second):
		t.Error("Test timed out - likely a Deadlock occurred")
	}
}

func TestHub_GracefulShutdown(t *testing.T) {
	h := NewHub()
	go h.Listen()

	// 1. Connect a client
	client := SseClient{
		Id:         "stuck-client",
		NotifyChan: make(chan database.Notification, 1),
	}
	h.AddClient(&client)

	// Wait for add
	time.Sleep(10 * time.Millisecond)

	// 2. Close the Hub
	// This should block until the client's channel is closed
	h.close()

	// 3. Verify the client was released
	select {
	case _, ok := <-client.NotifyChan:
		if ok {
			t.Error("Client channel received data instead of closing")
		}
		// Success: channel is closed (ok is false)
	case <-time.After(100 * time.Millisecond):
		t.Error("Client channel was NOT closed during shutdown (Memory Leak)")
	}
}
