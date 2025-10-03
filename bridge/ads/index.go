package ads

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/ads", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Ads management API service pinged")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})
}
