package config

import (
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watcher watches the config directory for changes
type Watcher struct {
	watcher  *fsnotify.Watcher
	dir      string
	onChange func()
	done     chan struct{}
}

// NewWatcher creates a new config directory watcher
func NewWatcher(dir string, onChange func()) (*Watcher, error) {
	w, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := w.Add(dir); err != nil {
		w.Close()
		return nil, err
	}

	return &Watcher{
		watcher:  w,
		dir:      dir,
		onChange: onChange,
		done:     make(chan struct{}),
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
				// Debounce: reset timer on each event
				if debounceTimer != nil {
					debounceTimer.Stop()
				}
				debounceTimer = time.AfterFunc(debounceDelay, func() {
					log.Printf("Config changed: %s (%s)", event.Name, event.Op)
					w.onChange()
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
