package api

import (
	"encoding/json"
	"net/http"
	"service/access"
	"service/database"
	"service/log"
	"service/utils"
)

func init() {
	http.HandleFunc("/api/report", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Receiving report...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

			var body struct {
				AdID        int64  `json:"ad_id"`
				AccountID   int    `json:"account_id"`
				AuthToken   string `json:"authtoken"`
				Description string `json:"description"`
			}

			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				log.Error("Failed to parse JSON body: %s", err.Error())
				http.Error(w, "Failed to parse JSON body", http.StatusInternalServerError)
				return
			}

			user := &utils.ArgonUser{Account: body.AccountID, Token: body.AuthToken}
			valid, err := access.ValidateArgonUser(user)
			if err != nil {
				log.Error("Failed to validate Argon user: %s", err.Error())
				http.Error(w, "Failed to validate Argon user", http.StatusInternalServerError)
				return
			}

			if valid {
				user, err = access.GetArgonUser(body.AccountID)
				if err != nil {
					log.Error("Failed to check for Argon user: %s", err.Error())
					http.Error(w, "Failed to check for Argon user", http.StatusInternalServerError)
					return
				}

				if user.ReportBanned {
					log.Warn("Argon user %s attempted to report ad of ID %d while banned", user.Account, body.AdID)
					http.Error(w, "Banned from ad reporting", http.StatusForbidden)
					return
				}

				err = database.NewReport(body.AdID, body.AccountID, body.Description)
				if err != nil {
					log.Error("Failed to create report: %s", err.Error())
					http.Error(w, "Failed to create report", http.StatusInternalServerError)
					return
				}

				log.Info("Registered report for ad of ID %d", body.AdID)
			} else {
				http.Error(w, "Invalid Argon user", http.StatusUnauthorized)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
