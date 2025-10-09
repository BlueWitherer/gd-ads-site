package stats

import (
	"net/http"

	"service/log"
)

func init() {
	http.HandleFunc("/stats", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Statistics API service pinged")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})
}
