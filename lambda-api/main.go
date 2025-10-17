package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// Function to gracefully shut down the server
func shutdown(server *http.Server) {
	log.Println("Shutting down the server...")

	// Create a deadline to wait for ongoing requests
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server Shutdown Failed:%+v", err)
	}
	log.Println("Server gracefully stopped")

	// Exit the app
	os.Exit(0)
}

func main() {
	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Handle /events/new endpoint
	mux.HandleFunc("/events/new", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received /events/new request")

		// Simulate some processing here if needed

		// Respond with 204 No Content
		w.WriteHeader(http.StatusNoContent)

		// Trigger shutdown in a separate goroutine
		go shutdown(server)
	})

	// Listen for system interrupt (optional in most cloud setups)
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		<-c
		shutdown(server)
	}()

	log.Println("Starting server on :8080")
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("Server error: %v", err)
	}
}
