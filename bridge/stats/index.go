package stats

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/stats", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Statistics API service pinged")
		w.Write([]byte("pong!"))
	})
}
