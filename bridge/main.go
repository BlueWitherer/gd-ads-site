package main

import (
	"net/http"
	"os"

	_ "bridge/ads"
	"bridge/log"
	_ "bridge/proxy"
	_ "bridge/stats"
)

func main() {
	log.Info("Starting server on http://localhost:8081")

	log.Debug("Serving static files")
	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)

	log.Debug("Starting handlers")
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server pinged!")
		asciiArt, err := os.ReadFile("../src/assets/aw-ascii.txt")
		if err != nil {
			log.Error("Failed to read ASCII art: " + err.Error())
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("pong!"))
			return
		}

		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write(asciiArt)
	})

	log.Done("Server started successfully")
	http.ListenAndServe(":8081", nil)

	log.Warn("Server stopped")
}
