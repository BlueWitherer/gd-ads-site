package ads

import (
	"encoding/json"
	"net/http"

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

			ads, err := database.FilterAdsByUser(rows, userID)
			if err != nil {
				log.Error("List ads failed: %s", err.Error())
				http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
				return
			}

			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(ads)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
