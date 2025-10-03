package main

import (
	"net/http"

	"bridge/log"
)

func main() {
	log.Info("Starting server on http://localhost:8080")

	http.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		log.Info("Server pinged!")
		w.Write([]byte("pong!"))
	})

	http.ListenAndServe(":8080", nil)
}
