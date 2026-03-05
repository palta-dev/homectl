package config

import (
	"context"
	"log"
	"time"

	"github.com/fsnotify/fsnotify"
)

// Watch monitors the config file for changes and calls the callback on updates
func Watch(ctx context.Context, path string, callback func(*Config)) error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	go func() {
		defer watcher.Close()

		var lastMod time.Time
		const debounce = 500 * time.Millisecond

		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				// Only process write events
				if event.Op&fsnotify.Write != fsnotify.Write {
					continue
				}

				// Debounce rapid changes
				now := time.Now()
				if now.Sub(lastMod) < debounce {
					continue
				}
				lastMod = now

				// Small delay to ensure file is fully written
				time.Sleep(100 * time.Millisecond)

				// Load new config
				newCfg, err := Load(path)
				if err != nil {
					log.Printf("Config reload failed: %v", err)
					continue
				}

				callback(newCfg)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Printf("Config watcher error: %v", err)
			}
		}
	}()

	if err := watcher.Add(path); err != nil {
		return err
	}

	log.Printf("Watching config file for changes: %s", path)
	return nil
}
