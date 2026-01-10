package server

import (
	"sync"
)

// Broadcaster manages SSE client connections and broadcasts messages
type Broadcaster struct {
	clients map[chan []byte]bool
	mu      sync.RWMutex
}

// NewBroadcaster creates a new broadcaster
func NewBroadcaster() *Broadcaster {
	return &Broadcaster{
		clients: make(map[chan []byte]bool),
	}
}

// Subscribe adds a new client and returns a channel for receiving messages
func (b *Broadcaster) Subscribe() chan []byte {
	ch := make(chan []byte, 16) // Buffer to prevent blocking
	b.mu.Lock()
	b.clients[ch] = true
	b.mu.Unlock()
	return ch
}

// Unsubscribe removes a client
func (b *Broadcaster) Unsubscribe(ch chan []byte) {
	b.mu.Lock()
	delete(b.clients, ch)
	close(ch)
	b.mu.Unlock()
}

// Broadcast sends data to all connected clients
func (b *Broadcaster) Broadcast(data []byte) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for ch := range b.clients {
		select {
		case ch <- data:
		default:
			// Client buffer full, skip this message
		}
	}
}

// ClientCount returns the number of connected clients
func (b *Broadcaster) ClientCount() int {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return len(b.clients)
}
