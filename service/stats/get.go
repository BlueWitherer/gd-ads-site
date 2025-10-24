package stats

import (
	"encoding/json"
	"net/http"

	"service/access"
	"service/database"
	"service/log"
)

type Stats struct {
	Views  int `json:"views"`
	Clicks int `json:"clicks"`
}

func init() {
	http.HandleFunc("/stats/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting advertisement stats for user...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Require logged-in user
		uid, err := access.GetSessionUserID(r)
		if err != nil || uid == "" {
			log.Error("Unauthorized access to /stats/get: %s", err.Error())
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		header.Set("Content-Type", "application/json")

		views, clicks, err := database.GetUserTotals(uid)
		if err != nil {
			log.Error("Failed to fetch user totals: %s", err.Error())
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}

		stats := Stats{
			Views:  views,
			Clicks: clicks,
		}

		log.Info("Retrieved stats for user: %s", uid)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error("Failed to encode stats response: %s", err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
