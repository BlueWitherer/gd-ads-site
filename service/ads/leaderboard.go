package ads

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/leaderboard", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Ads leaderboard API service pinged")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")
		header.Set("Content-Type", "text/plain")

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong!"))
	})

	http.HandleFunc("/ads/leaderboard/views", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			query := r.URL.Query()
			startStr := query.Get("start")
			endStr := query.Get("end")

			start, err := strconv.ParseUint(startStr, 10, 64)
			if err != nil {
				log.Error("Failed to get starting position: %s", err.Error())
				http.Error(w, "Failed to get starting position", http.StatusBadRequest)
				return
			}

			end, err := strconv.ParseUint(endStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ending position: %s", err.Error())
				http.Error(w, "Failed to get ending position", http.StatusBadRequest)
				return
			}

			users, err := database.UserLeaderboard(database.StatByViews, start, end)
			if err != nil {
				log.Error("Failed to get leaderboard: %s", err.Error())
				http.Error(w, "Failed to get leaderboard", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(users); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ads/leaderboard/clicks", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			query := r.URL.Query()
			startStr := query.Get("start")
			endStr := query.Get("end")

			start, err := strconv.ParseUint(startStr, 10, 64)
			if err != nil {
				log.Error("Failed to get starting position: %s", err.Error())
				http.Error(w, "Failed to get starting position", http.StatusBadRequest)
				return
			}

			end, err := strconv.ParseUint(endStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ending position: %s", err.Error())
				http.Error(w, "Failed to get ending position", http.StatusBadRequest)
				return
			}

			users, err := database.UserLeaderboard(database.StatByClicks, start, end)
			if err != nil {
				log.Error("Failed to get leaderboard: %s", err.Error())
				http.Error(w, "Failed to get leaderboard", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(users); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
