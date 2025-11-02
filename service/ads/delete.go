package ads

import (
	"net/http"
	"strconv"

	"service/access"
	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/delete", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Attempting to delete ad(s)...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "DELETE")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodDelete {
			header.Set("Content-Type", "application/json")

			uid, err := access.GetSessionUserID(r)
			if err != nil {
				log.Error("Unauthorized access to /ads/delete: %s", err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			idStr := r.URL.Query().Get("id")
			if idStr == "" {
				http.Error(w, "Missing ad ID parameter", http.StatusBadRequest)
				return
			}

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid ad ID parameter", http.StatusBadRequest)
				return
			}

			permission := false // does the user meet the criteria

			ownerid, err := database.GetAdvertisementOwnerId(id)
			if err != nil {
				log.Error("Failed to get advertisement owner: %s", err.Error())
				http.Error(w, "Failed to get advertisement owner", http.StatusInternalServerError)
				return
			}

			user, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user:", http.StatusInternalServerError)
				return
			}

			if user.IsAdmin || user.IsStaff || ownerid == user.ID {
				permission = true
			}

			if permission {
				ad, err := database.DeleteAdvertisement(id)
				if err != nil {
					log.Error("Failed to delete advertisement: %s", err.Error())
					http.Error(w, "Failed to delete advertisement", http.StatusInternalServerError)
					return
				}

				log.Info("Deleted advertisement of ID %d", ad.AdID)

				w.WriteHeader(http.StatusOK)
				w.Write([]byte(`{"status":"success","message":"Advertisement deleted successfully"}`))
			} else {
				log.Error("Unauthorized deletion attempt for ad ID %d by user %s", id, uid)
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
