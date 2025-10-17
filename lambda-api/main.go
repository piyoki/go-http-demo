// main.go
package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	// Use an atomic flag to ensure we only trigger shutdown once.
	var shuttingDown int32

	mux := http.NewServeMux()
	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Health endpoint (optional but handy for probes)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// /events/new endpoint: stateless processing, return 204, then stop the server
	mux.HandleFunc("/events/new", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// TODO: stateless processing here (validate, publish to queue, etc.)
		// e.g., read minimal data without storing local state.

		// Respond 204 (No Content)
		w.WriteHeader(http.StatusNoContent)

		// Trigger a graceful shutdown *once* after the response is sent.
		if atomic.CompareAndSwapInt32(&shuttingDown, 0, 1) {
			go func() {
				// Small delay helps ensure the 204 is flushed before Shutdown starts.
				time.Sleep(50 * time.Millisecond)

				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				if err := server.Shutdown(ctx); err != nil {
					log.Printf("graceful shutdown error: %v", err)
				}
			}()
		}
	})

	// Handle SIGTERM/SIGINT for graceful shutdown (K8s/containers)
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		s := <-sigCh
		log.Printf("received signal: %s; shutting down...", s)

		if atomic.CompareAndSwapInt32(&shuttingDown, 0, 1) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := server.Shutdown(ctx); err != nil {
				log.Printf("graceful shutdown error: %v", err)
			}
		}
	}()

	log.Printf("listening on %s", server.Addr)
	// ListenAndServe returns http.ErrServerClosed on graceful shutdown.
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}

	log.Println("server stopped cleanly")
}
