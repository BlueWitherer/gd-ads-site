package stats

import (
	"encoding/json"
	"net/http"

	"service/access"
	"service/database"
	"service/log"
)

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
		if err != nil {
			log.Error("Failed to get session ID: %s", err.Error())
			http.Error(w, "Failed to get session ID", http.StatusInternalServerError)
			return
		} else if uid == "" {
			log.Error("Unauthorized access to /stats/get")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		header.Set("Content-Type", "application/json")

		stats, err := database.GetUserTotals(uid)
		if err != nil {
			log.Error("Failed to fetch user totals: %s", err.Error())
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}

		log.Info("Retrieved stats for user: %s", uid)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error("Failed to encode stats response: %s", err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})

	http.HandleFunc("/stats/global", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting global advertisement statistics...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		// Get global stats from database
		stats, err := database.GetGlobalStats()
		if err != nil {
			log.Error("Failed to fetch global stats: %s", err.Error())
			http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
			return
		}

		log.Debug("Retrieved global stats - Views: %d, Clicks: %d, Ads: %d", stats.TotalViews, stats.TotalClicks, stats.AdCount)
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(stats); err != nil {
			log.Error("Failed to encode global stats response: %s", err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
