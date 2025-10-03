package main

import (
	"net/http"

	"bridge/log"
)

func main() {
	log.Info("Starting server on http://localhost:8080")

	log.Debug("Serving static files")
	fs := http.FileServer(http.Dir("../dist"))
	http.Handle("/", fs)

	log.Debug("Starting handlers")
	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server pinged!")
		w.Write([]byte("pong!"))
	})

	log.Done("Server started successfully")
	http.ListenAndServe(":8080", nil)

	log.Warn("Server stopped")
}
