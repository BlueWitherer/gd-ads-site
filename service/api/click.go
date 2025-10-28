package api

import (
	"encoding/json"
	"net/http"

	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/api/click", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Registering click...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			// Validate user agent is from the game client
			userAgent := r.Header.Get("User-Agent")
			if userAgent != "PlayerAdvertisements/1.0" {
				log.Warn("Click rejected: invalid user agent '%s'", userAgent)
				http.Error(w, "Unauthorized user agent", http.StatusForbidden)
				return
			}

			var body struct {
				AdID   int64  `json:"ad_id"`
				UserID string `json:"user_id"`
			}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				log.Error("Failed to parse JSON body: %s", err.Error())
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}

			// Validate user_id is not empty
			if body.UserID == "" {
				log.Error("User ID is empty")
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			// Use user_id string directly (it stays as a string for the database)
			err := database.NewStatWithUserID(database.AdEventClick, body.AdID, body.UserID)
			if err != nil {
				log.Error("Failed to create database click statistic: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("Click registered!"))
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
