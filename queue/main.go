package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/nicograef/go-playground/queue/api"
	"github.com/nicograef/go-playground/queue/queue"
)

func main() {
	server := &http.Server{
		Addr: ":3000",
	}

	appQueue, err := queue.LoadQueueFromJsonFile()
	if err != nil {
		fmt.Println("No existing queue found, creating a new one.")
		appQueue = queue.New()
	} else {
		fmt.Println("Loaded existing queue from file.")
	}

	http.HandleFunc("/enqueue", api.NewEnqueueHandler(appQueue))
	http.HandleFunc("/peek", api.NewPeekHandler(appQueue))
	http.HandleFunc("/dequeue", api.NewDequeueHandler(appQueue))

	// Set up signal handling for graceful shutdown
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		fmt.Println("Signal received, persisting queue...")
		appQueue.PersistToJsonFile()
		os.Exit(0)
	}()

	err = server.ListenAndServe()
	if err != nil {
		fmt.Println("Server error:", err)
		appQueue.PersistToJsonFile()
		panic(err)
	}
}
