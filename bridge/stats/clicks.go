package stats

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/stats/clicks", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Clicks endpoint hit")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("42"))
	})
}
