package ads

import (
	"net/http"

	"bridge/log"
)

func init() {
	http.HandleFunc("/api/ads/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Attempting to delete ad(s)...")
		log.Warn("This feature has not been implemented yet!")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusNotImplemented)
		http.Error(w, "Not implemented", http.StatusNotImplemented)
	})
}
