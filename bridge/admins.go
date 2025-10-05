package main

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/admins", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Admins API service pinged")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})

	http.HandleFunc("/admins/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Admin database API service pinged")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})

	http.HandleFunc("/admins/all", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Admin full database API service pinged")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})
}
