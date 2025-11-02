package ads

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service/database"
	"service/log"
	"service/utils"
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
			pageStr := query.Get("page")
			maxStr := query.Get("max")

			page, err := strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				log.Error("Failed to get starting position: %s", err.Error())
				http.Error(w, "Failed to get starting position", http.StatusBadRequest)
				return
			}

			max, err := strconv.ParseUint(maxStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ending position: %s", err.Error())
				http.Error(w, "Failed to get ending position", http.StatusBadRequest)
				return
			}

			users, err := database.UserLeaderboard(utils.StatByViews, page, max)
			if err != nil {
				log.Error("Failed to get views leaderboard: %s", err.Error())
				http.Error(w, "Failed to get views leaderboard", http.StatusInternalServerError)
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
			pageStr := query.Get("page")
			maxStr := query.Get("max")

			page, err := strconv.ParseUint(pageStr, 10, 64)
			if err != nil {
				log.Error("Failed to get starting position: %s", err.Error())
				http.Error(w, "Failed to get starting position", http.StatusBadRequest)
				return
			}

			max, err := strconv.ParseUint(maxStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ending position: %s", err.Error())
				http.Error(w, "Failed to get ending position", http.StatusBadRequest)
				return
			}

			users, err := database.UserLeaderboard(utils.StatByClicks, page, max)
			if err != nil {
				log.Error("Failed to get clicks leaderboard: %s", err.Error())
				http.Error(w, "Failed to get clicks leaderboard", http.StatusInternalServerError)
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
