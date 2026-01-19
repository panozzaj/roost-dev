package config

import (
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches the config directory for changes
type Watcher struct {
	watcher  *fsnotify.Watcher
	dir      string
	onChange func(changedFiles []string)
	done     chan struct{}

	// Track changed files during debounce window
	pendingMu    sync.Mutex
	pendingFiles map[string]bool
}

// NewWatcher creates a new config directory watcher
// The onChange callback receives a list of changed filenames (base names, not full paths)
func NewWatcher(dir string, onChange func(changedFiles []string)) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := w.Add(dir); err != nil {
		w.Close()
		return nil, err
	}

	return &Watcher{
		watcher:      w,
		dir:          dir,
		onChange:     onChange,
		done:         make(chan struct{}),
		pendingFiles: make(map[string]bool),
	}, nil
}

// Start begins watching for changes
func (w *Watcher) Start() {
	go w.run()
}

// Stop stops the watcher
func (w *Watcher) Stop() {
	close(w.done)
	w.watcher.Close()
}

func (w *Watcher) run() {
	// Debounce timer - wait for rapid changes to settle
	var debounceTimer *time.Timer
	debounceDelay := 200 * time.Millisecond

	for {
		select {
		case <-w.done:
			if debounceTimer != nil {
				debounceTimer.Stop()
			}
			return

		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}

			// Only react to relevant events
			if event.Op&(fsnotify.Write|fsnotify.Create|fsnotify.Remove|fsnotify.Rename) != 0 {
				// Track this changed file
				w.pendingMu.Lock()
				w.pendingFiles[filepath.Base(event.Name)] = true
				w.pendingMu.Unlock()

				// Debounce: reset timer on each event
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(debounceDelay, func() {
					// Collect and clear pending files
					w.pendingMu.Lock()
					files := make([]string, 0, len(w.pendingFiles))
					for f := range w.pendingFiles {
						files = append(files, f)
					}
					w.pendingFiles = make(map[string]bool)
					w.pendingMu.Unlock()

					log.Printf("Config changed: %v", files)
					w.onChange(files)
				})
			}

		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			log.Printf("Config watcher error: %v", err)
		}
	}
}
