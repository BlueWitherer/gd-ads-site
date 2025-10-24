package ads

import (
	"encoding/json"
	"net/http"
	"strconv"

	"service/access"
	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/pending", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if !u.IsAdmin {
				log.Error("User of ID %s is not admin", u.ID)
				http.Error(w, "User is not admin", http.StatusUnauthorized)
				return
			}

			rows, err := database.ListAllAdvertisements()
			if err != nil {
				log.Error("Failed to list ads: %s", err.Error())
				http.Error(w, "Failed to list ads", http.StatusInternalServerError)
				return
			}

			var adList []database.Ad

			list, err := database.FilterAdsByPending(rows, true)
			if err != nil {
				log.Error("Failed to filter pending ads: %s", err.Error())
				http.Error(w, "Failed to filter pending ads", http.StatusInternalServerError)
				return
			}

			query := r.URL.Query()
			user := query.Get("user")

			if user != "" {
				list, err = database.FilterAdsByUser(adList, user)
				if err != nil {
					log.Error("List ads failed: %s", err.Error())
					http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
					return
				}
			}

			adList = append(adList, list...)

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(adList); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ads/pending/accept", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			w.Header().Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if !u.IsAdmin {
				log.Error("User of ID %s is not admin", u.ID)
				http.Error(w, "User is not admin", http.StatusUnauthorized)
				return
			}

			query := r.URL.Query()
			idStr := query.Get("id")

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ad ID: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ad, err := database.ApproveAd(id)
			if err != nil {
				log.Error("Failed to approve ad: %s", err.Error())
				http.Error(w, "Failed to approve ad", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ad); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
