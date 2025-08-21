package main

import (
	"github.com/fsnotify/fsnotify"
	"log"
)

func main() {
	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Has(fsnotify.Write) {
					log.Println("file event:", event.Name)
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("errors:", err)
			}
		}
	}()

	// Add a path.
	err = watcher.Add("C:\\Users\\Administrator\\Desktop\\111")
	if err != nil {
		log.Fatal(err)
	}

	err = watcher.Add("C:\\Users\\Administrator\\Desktop\\222")
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}
