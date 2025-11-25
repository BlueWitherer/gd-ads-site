package api

import (
	"encoding/json"
	"net/http"
	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/api/announcement", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting latest announcement...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			announcement, err := database.GetLatestAnnouncement()
			if err != nil {
				log.Error("Failed to get latest announcement: %s", err.Error())
				http.Error(w, "Failed to get latest announcement", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(announcement); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
