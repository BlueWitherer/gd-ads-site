package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/api/view", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Registering view...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			var body struct {
				AdID   int64  `json:"ad_id"`
				UserID string `json:"user_id"`
			}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				log.Error("Failed to parse JSON body: %s", err.Error())
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}

			user, err := strconv.ParseInt(body.UserID, 10, 64)
			if err != nil {
				log.Error("Failed to parse user ID: %s", err.Error())
				http.Error(w, "Invalid user ID", http.StatusBadRequest)
				return
			}

			err = database.NewStat(database.AdEventView, body.AdID, user)
			if err != nil {
				log.Error("Failed to create database view statistic: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
			} else {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("View registered!"))
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
