package stats

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Statistics API service pinged")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})
}
