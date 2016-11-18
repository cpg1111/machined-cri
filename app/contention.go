package app

import (
	"log"

	"golang.org/x/exp/inotify"
)

func watchForContention(path string, done chan struct{}) error {
	watcher, wErr := inotify.NewWatcher()
	if wErr != nil {
		log.Println(wErr.Error())
		return wErr
	}
	wErr = watcher.AddWatch(path, inotify.IN_OPEN|inotify.IN_DELETE_SELF)
	if wErr != nil {
		log.Println(wErr.Error())
		return wErr
	}
	go func() {
		select {
		case evt := <-watcher.Event:
			log.Printf("inotify event: %v", evt)
		case err := <-watcher.Error:
			log.Printf("inotify watcher error: %v", err)
		}
		close(done)
	}()
	return nil
}
