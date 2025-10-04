package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "bridge/ads"
	"bridge/log"
	_ "bridge/proxy"
	_ "bridge/stats"
)

func main() {
	log.Print("Starting server...")

	log.Debug("Serving static files")
	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)

	log.Debug("Starting handlers")
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server pinged!")
		asciiArt, err := os.ReadFile("../src/assets/aw-ascii.txt")
		if err != nil {
			log.Error("Failed to read ASCII art: " + err.Error())
			header := w.Header()

			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET")
			header.Set("Access-Control-Allow-Headers", "Content-Type")
			header.Set("Content-Type", "text/plain")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
		} else {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write(asciiArt)
		}
	})

	srv := &http.Server{Addr: ":8081"}

	go func() {
		log.Done("Server started successfully on host http://localhost:8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Warn("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Shutdown error: " + err.Error())
	} else {
		log.Print("Server stopped")
	}
}
