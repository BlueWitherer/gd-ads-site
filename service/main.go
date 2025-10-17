package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"service/access"
	_ "service/ads"
	_ "service/database"
	"service/log"
	_ "service/proxy"
	_ "service/stats"
)

func expiryCleanupRoutine(adFolder string) {
	go func() {
		for {
			log.Info("Scanning for expired %s ads...", adFolder)

			adsDir := filepath.Join("..", "ad_storage", adFolder)
			files, _ := os.ReadDir(adsDir)
			for _, file := range files {
				info, err := file.Info()
				if err != nil {
					log.Error("Failed to get ad file info for %s: %s", file.Name(), err.Error())
					continue
				}

				if time.Since(info.ModTime()) > 7*24*time.Hour {
					os.Remove(filepath.Join(adsDir, file.Name()))
				}
			}

			time.Sleep(12 * time.Hour) // Run twice a day
		}
	}()
}

func main() {
	log.Print("Starting server...")

	log.Debug("Starting handlers...")
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server pinged!")
		asciiArt, err := os.ReadFile("../src/assets/aw-ascii.txt")
		if err != nil {
			log.Error("Failed to read ASCII art: %s", err.Error())
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

	// SPA fallback
	log.Debug("Setting up SPA fallback for client-side routing")
	staticDir := "../dist"
	fs := http.FileServer(http.Dir(staticDir))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Received request for host %s", access.FullURL(r))

		requestedPath := strings.TrimPrefix(filepath.Clean(r.URL.Path), "/")
		fullPath := filepath.Join(staticDir, requestedPath)
		if requestedPath == "" || requestedPath == "." {
			http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
			return
		}

		info, err := os.Stat(fullPath)
		if err == nil && !info.IsDir() {
			fs.ServeHTTP(w, r)
			return
		}

		log.Debug("Serving index.html for SPA route: %s", r.URL.Path)
		http.ServeFile(w, r, filepath.Join(staticDir, "index.html"))
	})

	srv := &http.Server{Addr: ":8081"}

	go func() {
		log.Debug("Starting expiry routines...")
		expiryCleanupRoutine("banner")
		expiryCleanupRoutine("skyscraper")
		expiryCleanupRoutine("square")

		log.Done("Server started successfully on host http://localhost%s", srv.Addr)
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
