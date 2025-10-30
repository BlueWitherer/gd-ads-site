package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"service/access"
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
			query := r.URL.Query()
			accountIDStr := query.Get("account_id")
			authToken := query.Get("authtoken")

			log.Info("Received view request - account_id=%s, authtoken length=%d", accountIDStr, len(authToken))

			var body struct {
				AdID   int64  `json:"ad_id"`
				UserID string `json:"user_id"`
			}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				log.Error("Failed to parse JSON body: %s", err.Error())
				http.Error(w, "Invalid JSON body", http.StatusBadRequest)
				return
			}

			log.Debug("Body decoded - AdID: %v, UserID: %s", body.AdID, body.UserID)

			// Validate required parameters
			if accountIDStr == "" || authToken == "" {
				log.Error("Missing query parameters - account_id: %s, authtoken: %s", accountIDStr, authToken)
				http.Error(w, "Missing account_id or authtoken", http.StatusBadRequest)
				return
			}

			// Convert account_id to int
			var accountID int
			_, err := fmt.Sscanf(accountIDStr, "%d", &accountID)
			if err != nil {
				log.Error("Failed to parse account_id: %s, error: %s", accountIDStr, err.Error())
				http.Error(w, "Invalid account_id format", http.StatusBadRequest)
				return
			}

			user := access.ArgonUser{Account: accountID, Token: authToken}
			valid, err := access.ValidateArgonUser(user)
			if err != nil {
				log.Error("Failed to validate Argon user: %s", err.Error())
				http.Error(w, "Failed to validate Argon user", http.StatusUnauthorized)
				return
			}

			if valid {
				// Validate user_id is not empty
				if body.UserID == "" {
					log.Error("User ID is empty")
					http.Error(w, "Invalid user ID", http.StatusBadRequest)
					return
				}

				// Convert user_id string directly (it stays as a string for the database)
				err := database.NewStatWithUserID(database.AdEventView, body.AdID, body.UserID)
				if err != nil {
					log.Error("Failed to create database view statistic: %s", err.Error())
					http.Error(w, err.Error(), http.StatusInternalServerError)
				} else {
					w.WriteHeader(http.StatusOK)
					w.Write([]byte("View registered!"))
				}
			} else {
				http.Error(w, "Argon user invalid", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
