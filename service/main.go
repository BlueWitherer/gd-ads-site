package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"sync"
	"syscall"
	"time"

	"service/access"
	_ "service/ads"
	_ "service/api"
	"service/database"
	"service/log"
	_ "service/proxy"
	_ "service/stats"

	"golang.org/x/time/rate"
)

var visitors = make(map[string]*rate.Limiter)
var mu sync.Mutex
var visitorCleanupTicker *time.Ticker

const VISITOR_EXPIRY = 6 * time.Hour

func getClientIP(r *http.Request) string {
	if cf := r.Header.Get("CF-Connecting-IP"); cf != "" {
		return cf
	}
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		return strings.Split(xff, ",")[0]
	}
	return strings.Split(r.RemoteAddr, ":")[0]
}

func getVisitor(ip string) *rate.Limiter {
	mu.Lock()
	defer mu.Unlock()

	limiter, exists := visitors[ip]
	if !exists {
		limiter = rate.NewLimiter(15, 40) // generous: 10 req/sec, burst of 20
		visitors[ip] = limiter
	}

	return limiter
}

// Starts a background goroutine to clean up rate limiters periodically
func startVisitorCleanup() {
	visitorCleanupTicker = time.NewTicker(6 * time.Hour)
	go func() {
		for range visitorCleanupTicker.C {
			cleanOldVisitors()
		}
	}()
	log.Info("Visitor rate limiter cleanup goroutine started - will clean every 6 hours")
}

// Stops the visitor cleanup goroutine
func stopVisitorCleanup() {
	if visitorCleanupTicker != nil {
		visitorCleanupTicker.Stop()
		log.Info("Visitor rate limiter cleanup goroutine stopped")
	}
}

// Cleans up the visitors map to prevent unbounded growth
func cleanOldVisitors() {
	mu.Lock()
	defer mu.Unlock()

	// Limiting visitors map to reasonable size - keep it under 10k IPs
	// This prevents DDoS-like memory exhaustion from many unique IPs
	if len(visitors) > 10000 {
		// Clear all visitors and start fresh
		visitors = make(map[string]*rate.Limiter)
		log.Warn("Visitor rate limiter map cleared - exceeded 10k unique IPs")
	} else {
		log.Debug("Visitor rate limiter map size: %d entries", len(visitors))
	}
}

func rateLimitMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		limiter := getVisitor(ip)

		if !limiter.Allow() {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func expiryCleanupRoutine(adFolder string) {
	go func() {
		for {
			log.Debug("Scanning for expired %s ads...", adFolder)

			adsDir := filepath.Join("..", "ad_storage", adFolder)
			files, _ := os.ReadDir(adsDir)
			for _, file := range files {
				info, err := file.Info()
				if err != nil {
					log.Error("Failed to get ad file info for %s: %s", file.Name(), err.Error())
				} else {
					if time.Since(info.ModTime()) > 7*24*time.Hour {
						log.Info("Removing expired ad %s (%v B)", info.Name(), info.Size())
						os.Remove(filepath.Join(adsDir, file.Name()))
					} else {
						log.Debug("Ad %s is still valid", info.Name())
					}
				}
			}

			time.Sleep(12 * time.Hour) // Run twice a day
		}
	}()
}

func expiryCleanupRoutineSql() {
	go func() {
		for {
			log.Debug("Deleting expired ad records...")
			err := database.DeleteAllExpiredAds()
			if err != nil {
				log.Error("Failed to delete expired ad records: %s", err.Error())
			} else {
				log.Info("Expired ad records cleanup complete")
			}

			time.Sleep(12 * time.Hour) // Run twice a day
		}
	}()
}

func main() {
	log.Print("Starting server...")
	access.StartSessionCleanup()
	startVisitorCleanup()

	srv := &http.Server{Addr: fmt.Sprintf(":%s", os.Getenv("WEB_PORT"))}

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

	log.Debug("Starting image handler...")
	http.HandleFunc("/cdn/", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		requestedPath := strings.TrimPrefix(r.URL.Path, "/cdn/")
		fullPath := filepath.Join("../ad_storage", requestedPath)

		// Set cache control headers to prevent browser caching of ad images
		header.Set("Content-Type", "image/webp")
		header.Set("Cache-Control", "no-cache, no-store, must-revalidate")
		header.Set("Pragma", "no-cache")
		header.Set("Expires", "0")

		http.ServeFile(w, r, fullPath)
	})

	log.Debug("Starting handlers...")
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server pinged!")
		header := w.Header()

		asciiArt, err := os.ReadFile("../src/assets/aw-ascii.txt")
		if err != nil {
			log.Error("Failed to read ASCII art: %s", err.Error())

			header.Set("Access-Control-Allow-Origin", "*")
			header.Set("Access-Control-Allow-Methods", "GET")
			header.Set("Access-Control-Allow-Headers", "Content-Type")
			header.Set("Content-Type", "text/plain")

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
		} else {
			header.Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write(asciiArt)
		}
	})

	go func() {
		log.Debug("Starting expiry routines...")
		expiryCleanupRoutine("banner")
		expiryCleanupRoutine("skyscraper")
		expiryCleanupRoutine("square")

		log.Debug("Starting expiry routines for database side...")
		expiryCleanupRoutineSql()

		log.Done("Server started successfully on host http://localhost%s", srv.Addr)
		srv.Handler = rateLimitMiddleware(http.DefaultServeMux)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error(err.Error())
		}
	}()

	// Wait for interrupt signal to gracefully shut down
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Warn("Shutting down server...")

	// Stop all cleanup goroutines
	access.StopSessionCleanup()
	stopVisitorCleanup()

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Shutdown error: %s", err.Error())
	} else {
		log.Print("Server stopped")
	}
}
