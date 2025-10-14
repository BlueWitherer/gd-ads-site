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
		// require login
		userID, err := access.GetSessionUserID(r)
		if err != nil || userID == "" {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		ads, err := database.ListAdvertisementsByUser(userID)
		if err != nil {
			log.Error("List ads failed: %s", err.Error())
			http.Error(w, "Failed to fetch ads", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(ads)
	})
}
