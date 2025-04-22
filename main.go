// main.go
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/itsyaboikris/go_socket_server/server"
)

func main() {
	// Create a context that will be canceled on interrupt
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		log.Println("Starting server...")
		server.Start(ctx, ":9000")
	}()

	// Wait for interrupt signal
	<-sigChan
	log.Println("\nShutdown signal received")

	// Cancel context to initiate graceful shutdown
	cancel()
	log.Println("Server shutdown complete")
}
