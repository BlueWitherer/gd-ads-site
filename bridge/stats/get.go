package stats

import (
	"encoding/json"
	"net/http"

	"bridge/log"
)

type Stats struct {
	Views  int `json:"views"`
	Clicks int `json:"clicks"`
}

func init() {
	http.HandleFunc("/api/stats/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting advertisement stats for user...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "application/json")

		log.Debug("Constructing stats object response...")
		stats := Stats{
			Views:  906,
			Clicks: 24,
		}

		log.Info("Retrieved stats for user")
		w.WriteHeader(http.StatusOK)
		err := json.NewEncoder(w).Encode(stats)
		if err != nil {
			log.Error("Failed to encode stats response: " + err.Error())
			http.Error(w, "Failed to encode response", http.StatusInternalServerError)
			return
		}
	})
}
