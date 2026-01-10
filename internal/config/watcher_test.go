package config

import (
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"
)

func TestWatcher(t *testing.T) {
	t.Run("detects file changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yml")

		// Create initial file
		if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		var changed atomic.Bool
		w, err := NewWatcher(tmpDir, func() {
			changed.Store(true)
		})
		if err != nil {
			t.Fatalf("failed to create watcher: %v", err)
		}
		w.Start()
		defer w.Stop()

		// Give the watcher time to start
		time.Sleep(50 * time.Millisecond)

		// Modify the file
		if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
			t.Fatalf("failed to modify test file: %v", err)
		}

		// Wait for debounce + processing
		time.Sleep(400 * time.Millisecond)

		if !changed.Load() {
			t.Error("expected onChange to be called after file modification")
		}
	})

	t.Run("detects new files", func(t *testing.T) {
		tmpDir := t.TempDir()

		var changed atomic.Bool
		w, err := NewWatcher(tmpDir, func() {
			changed.Store(true)
		})
		if err != nil {
			t.Fatalf("failed to create watcher: %v", err)
		}
		w.Start()
		defer w.Stop()

		// Give the watcher time to start
		time.Sleep(50 * time.Millisecond)

		// Create a new file
		testFile := filepath.Join(tmpDir, "new.yml")
		if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		// Wait for debounce + processing
		time.Sleep(400 * time.Millisecond)

		if !changed.Load() {
			t.Error("expected onChange to be called after new file creation")
		}
	})

	t.Run("detects file deletion", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yml")

		// Create initial file
		if err := os.WriteFile(testFile, []byte("content"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		var changed atomic.Bool
		w, err := NewWatcher(tmpDir, func() {
			changed.Store(true)
		})
		if err != nil {
			t.Fatalf("failed to create watcher: %v", err)
		}
		w.Start()
		defer w.Stop()

		// Give the watcher time to start
		time.Sleep(50 * time.Millisecond)

		// Delete the file
		if err := os.Remove(testFile); err != nil {
			t.Fatalf("failed to delete test file: %v", err)
		}

		// Wait for debounce + processing
		time.Sleep(400 * time.Millisecond)

		if !changed.Load() {
			t.Error("expected onChange to be called after file deletion")
		}
	})

	t.Run("debounces rapid changes", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yml")

		// Create initial file
		if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		var callCount atomic.Int32
		w, err := NewWatcher(tmpDir, func() {
			callCount.Add(1)
		})
		if err != nil {
			t.Fatalf("failed to create watcher: %v", err)
		}
		w.Start()
		defer w.Stop()

		// Give the watcher time to start
		time.Sleep(50 * time.Millisecond)

		// Make rapid changes
		for i := 0; i < 5; i++ {
			if err := os.WriteFile(testFile, []byte("change "+string(rune('0'+i))), 0644); err != nil {
				t.Fatalf("failed to modify test file: %v", err)
			}
			time.Sleep(50 * time.Millisecond) // Less than debounce delay
		}

		// Wait for debounce to settle
		time.Sleep(400 * time.Millisecond)

		// Should only have been called once due to debouncing
		count := callCount.Load()
		if count != 1 {
			t.Errorf("expected onChange to be called 1 time (debounced), got %d", count)
		}
	})

	t.Run("Stop prevents further callbacks", func(t *testing.T) {
		tmpDir := t.TempDir()
		testFile := filepath.Join(tmpDir, "test.yml")

		// Create initial file
		if err := os.WriteFile(testFile, []byte("initial"), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}

		var called atomic.Bool
		w, err := NewWatcher(tmpDir, func() {
			called.Store(true)
		})
		if err != nil {
			t.Fatalf("failed to create watcher: %v", err)
		}
		w.Start()

		// Stop immediately
		w.Stop()

		// Make a change after stopping
		if err := os.WriteFile(testFile, []byte("modified"), 0644); err != nil {
			t.Fatalf("failed to modify test file: %v", err)
		}

		// Wait longer than debounce
		time.Sleep(400 * time.Millisecond)

		if called.Load() {
			t.Error("onChange should not be called after Stop")
		}
	})

	t.Run("handles non-existent directory", func(t *testing.T) {
		_, err := NewWatcher("/nonexistent/path/12345", func() {})
		if err == nil {
			t.Error("expected error for non-existent directory")
		}
	})
}
