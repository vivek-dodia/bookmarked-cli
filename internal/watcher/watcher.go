package watcher

import (
	"fmt"
	"log"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

type Watcher struct {
	watcher     *fsnotify.Watcher
	debounceMs  int
	onChange    func()
	debounceTimer *time.Timer
}

// New creates a new file watcher with debouncing
func New(debounceMs int, onChange func()) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("failed to create watcher: %w", err)
	}

	return &Watcher{
		watcher:    watcher,
		debounceMs: debounceMs,
		onChange:   onChange,
	}, nil
}

// Watch starts watching the specified file
func (w *Watcher) Watch(filePath string) error {
	// Watch the parent directory since Chrome does atomic writes (creates temp, then renames)
	dir := filepath.Dir(filePath)
	fileName := filepath.Base(filePath)

	if err := w.watcher.Add(dir); err != nil {
		return fmt.Errorf("failed to watch directory: %w", err)
	}

	log.Printf("Watching for changes to: %s", filePath)

	go func() {
		for {
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					return
				}

				// Only process events for our target file
				if filepath.Base(event.Name) != fileName {
					continue
				}

				// Handle Write and Rename events (Chrome does atomic writes via rename)
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					log.Printf("Detected change: %s", event.Op.String())
					w.debounce()
				}

			case err, ok := <-w.watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Watcher error: %v", err)
			}
		}
	}()

	return nil
}

// debounce delays the onChange callback to avoid excessive calls
func (w *Watcher) debounce() {
	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}

	w.debounceTimer = time.AfterFunc(time.Duration(w.debounceMs)*time.Millisecond, func() {
		log.Println("Debounce period elapsed, triggering sync...")
		w.onChange()
	})
}

// Close stops the watcher
func (w *Watcher) Close() error {
	if w.debounceTimer != nil {
		w.debounceTimer.Stop()
	}
	return w.watcher.Close()
}
