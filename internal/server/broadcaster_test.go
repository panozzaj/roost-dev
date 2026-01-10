package server

import (
	"testing"
	"time"
)

func TestBroadcaster(t *testing.T) {
	t.Run("Subscribe and Unsubscribe", func(t *testing.T) {
		b := NewBroadcaster()

		// Initially no clients
		if b.ClientCount() != 0 {
			t.Errorf("expected 0 clients, got %d", b.ClientCount())
		}

		// Subscribe
		ch := b.Subscribe()
		if b.ClientCount() != 1 {
			t.Errorf("expected 1 client, got %d", b.ClientCount())
		}

		// Unsubscribe
		b.Unsubscribe(ch)
		if b.ClientCount() != 0 {
			t.Errorf("expected 0 clients after unsubscribe, got %d", b.ClientCount())
		}
	})

	t.Run("Broadcast sends to all clients", func(t *testing.T) {
		b := NewBroadcaster()

		ch1 := b.Subscribe()
		ch2 := b.Subscribe()
		defer b.Unsubscribe(ch1)
		defer b.Unsubscribe(ch2)

		// Broadcast a message
		msg := []byte("test message")
		b.Broadcast(msg)

		// Both channels should receive the message
		select {
		case received := <-ch1:
			if string(received) != string(msg) {
				t.Errorf("ch1: expected %q, got %q", msg, received)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("ch1: timeout waiting for message")
		}

		select {
		case received := <-ch2:
			if string(received) != string(msg) {
				t.Errorf("ch2: expected %q, got %q", msg, received)
			}
		case <-time.After(100 * time.Millisecond):
			t.Error("ch2: timeout waiting for message")
		}
	})

	t.Run("Broadcast skips full client buffers", func(t *testing.T) {
		b := NewBroadcaster()
		ch := b.Subscribe()
		defer b.Unsubscribe(ch)

		// Fill the buffer (buffer size is 16)
		for i := 0; i < 20; i++ {
			b.Broadcast([]byte("message"))
		}

		// Should not panic or block, and channel should have 16 messages
		count := 0
		for {
			select {
			case <-ch:
				count++
			default:
				goto done
			}
		}
	done:
		if count != 16 {
			t.Errorf("expected 16 buffered messages, got %d", count)
		}
	})

	t.Run("Multiple subscribers", func(t *testing.T) {
		b := NewBroadcaster()

		// Subscribe multiple clients
		channels := make([]chan []byte, 5)
		for i := range channels {
			channels[i] = b.Subscribe()
		}

		if b.ClientCount() != 5 {
			t.Errorf("expected 5 clients, got %d", b.ClientCount())
		}

		// Unsubscribe some
		b.Unsubscribe(channels[0])
		b.Unsubscribe(channels[2])

		if b.ClientCount() != 3 {
			t.Errorf("expected 3 clients after unsubscribe, got %d", b.ClientCount())
		}

		// Clean up remaining
		b.Unsubscribe(channels[1])
		b.Unsubscribe(channels[3])
		b.Unsubscribe(channels[4])
	})

	t.Run("Concurrent subscribe and broadcast", func(t *testing.T) {
		b := NewBroadcaster()
		done := make(chan struct{})

		// Spawn subscribers
		go func() {
			for i := 0; i < 10; i++ {
				ch := b.Subscribe()
				go func(c chan []byte) {
					<-c // Read at least one message
					b.Unsubscribe(c)
				}(ch)
			}
			close(done)
		}()

		// Broadcast concurrently
		for i := 0; i < 20; i++ {
			b.Broadcast([]byte("test"))
		}

		<-done
		// Test passes if no race conditions or deadlocks
	})
}
