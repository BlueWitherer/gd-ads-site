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
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Check if status=pending query parameter is present
			pending := r.URL.Query().Get("pending")
			if pending == "1" {
				// Get user to check if admin
				user, err := database.GetUser(uid)
				if err != nil {
					log.Error("Failed to get user: %s", err.Error())
					http.Error(w, "Failed to get user", http.StatusInternalServerError)
					return
				}

				// Only admins can view all pending ads
				if !user.IsAdmin {
					http.Error(w, "Unauthorized", http.StatusUnauthorized)
					return
				}

				// Get pending ads directly from database
				pendingAds, err := database.ListPendingAdvertisements()
				if err != nil {
					log.Error("Failed to list pending ads: %s", err.Error())
					http.Error(w, "Failed to list pending ads", http.StatusInternalServerError)
					return
				}

				log.Info("Found %d pending advertisements", len(pendingAds))
				for i, ad := range pendingAds {
					log.Info("Pending ad %d: ID=%d, UserID=%s, Pending=%v", i, ad.AdID, ad.UserID, ad.Pending)
				}

				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(pendingAds); err != nil {
					log.Error("Failed to encode response: %s", err.Error())
					http.Error(w, "Failed to encode response", http.StatusInternalServerError)
					return
				}

				return
			}

			// Default behavior: get user's own ads
			rows, err := database.ListAllAdvertisements()
			if err != nil {
				log.Error("Failed to list ads: %s", err.Error())
				http.Error(w, "Failed to list ads", http.StatusInternalServerError)
				return
			}

			filtered, err := database.FilterAdsByUser(rows, uid)
			if err != nil {
				log.Error("List ads failed: %s", err.Error())
				http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(filtered); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
