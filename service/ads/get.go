package ads

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"service/access"
	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/get", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			// require login
			userID, err := access.GetSessionUserID(r)
			if err != nil || userID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			rows, err := database.ListAllAdvertisements()
			if err != nil {
				log.Error("Failed to list ads: %s", err.Error())
				http.Error(w, "Failed to list ads", http.StatusInternalServerError)
				return
			}

			filtered, err := database.FilterAdsByUser(rows, userID)
			if err != nil {
				log.Error("List ads failed: %s", err.Error())
				http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
				return
			}

			// Map to client-friendly shape
			type OutAd struct {
				ID         int64  `json:"id"`
				Type       string `json:"type"`
				LevelID    string `json:"levelId"`
				Image      string `json:"image"`
				Expiration string `json:"expiration"`
			}

			var out []OutAd
			for _, a := range filtered {
				// map numeric type to string
				t := "unknown"
				switch a.Type {
				case 1:
					t = "Banner"
				case 2:
					t = "Square"
				case 3:
					t = "Skyscraper"
				}

				// expiration is computed as days until 7 days after created_at
				expiration := ""
				if parsed, err := time.Parse(time.RFC3339, a.Created); err == nil {
					daysLeft := 7 - int(time.Since(parsed).Hours()/24)
					if daysLeft > 1 {
						expiration = fmt.Sprintf("%d days left", daysLeft)
					} else if daysLeft == 1 {
						expiration = "1 day left"
					} else if daysLeft <= 0 {
						expiration = "expired"
					}
				}

				out = append(out, OutAd{
					ID:         a.AdID,
					Type:       t,
					LevelID:    a.LevelID,
					Image:      a.ImageURL,
					Expiration: expiration,
				})
			}

			w.Header().Set("Content-Type", "application/json")
			if err := json.NewEncoder(w).Encode(out); err != nil {
				log.Error("Failed to encode ad response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
