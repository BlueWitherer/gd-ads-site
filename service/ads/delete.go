package ads

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
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

			ownerid, err := database.GetAdvertisementOwner(id)
			if err != nil {
				log.Error("Failed to get advertisement owner: %s", err.Error())
				http.Error(w, "Failed to get advertisement owner", http.StatusInternalServerError)
				return
			}

			uid, err := access.GetSessionUserID(r)
			if err != nil {
				log.Error("Unauthorized access to /ads/delete: %s", err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			if ownerid == uid {
				ad, err := database.DeleteAdvertisement(id)
				if err != nil {
					log.Error("Failed to delete advertisement: %s", err.Error())
					http.Error(w, "Failed to delete advertisement", http.StatusInternalServerError)
					return
				}

				adFolder, err := database.AdTypeFromInt(ad.Type)
				if err != nil {
					log.Error("Failed to get advertisement folder: %s", err.Error())
					http.Error(w, "Failed to get advertisement folder", http.StatusInternalServerError)
					return
				}

				target := filepath.Join("..", "ad_storage", string(adFolder), fmt.Sprintf("%s.webp", ad.UserID))

				err = os.Remove(target)
				if err != nil {
					log.Error("Failed to delete advertisement image: %s", err.Error())
					http.Error(w, "Failed to delete advertisement image", http.StatusInternalServerError)
					return
				}

				log.Info("Deleted advertisement of ID %s", idStr)

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
