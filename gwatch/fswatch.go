package gwatch

import (
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
)

func StartSession(path string, command string) {

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	var lastEvent fsnotify.Event

	timer := time.NewTimer(time.Millisecond)
	<-timer.C

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				lastEvent = event
				timer.Reset(time.Millisecond * 100)

			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("error:", err)
			case <-timer.C:
				if lastEvent.Has(fsnotify.Write) {
					// log.Println("modified file:", lastEvent.Name)
					runCmd(command)
				}
			}
		}
	}()

	// Add a path.
	err = watcher.Add(path)
	if err != nil {
		log.Fatal(err)
	}

	// Block main goroutine forever.
	<-make(chan struct{})
}

func runCmd(command string) {

	cmdParts := strings.Split(command, " ")
	c := exec.Command(cmdParts[0], cmdParts[1:]...)
	// c := exec.Command("Invoke-Expression", command)
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Run()
	// if err := c.Run(); err != nil {
	// 	log.Print(err)
	// }
}
