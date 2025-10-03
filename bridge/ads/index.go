package ads

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/ads", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Ads management API service pinged")
		w.Write([]byte("pong!"))
	})
}
