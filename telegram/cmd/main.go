package main

import (
	"os"
	"os/signal"
	"syscall"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)

	// Wait for interrupt signal
	<-done

	// Create a channel that will never receive data
	forever := make(chan struct{})

	// Block indefinitely
	<-forever
}
