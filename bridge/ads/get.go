package ads

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/ads/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Attempting to get ad(s)...")
		log.Warn("This feature has not been implemented yet!")

		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})
}
