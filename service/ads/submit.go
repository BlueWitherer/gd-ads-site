package ads

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"service/access"
	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/submit", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

			// Require logged-in session
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Parse form with 10MB limit
			r.ParseMultipartForm(10 << 20)

			// Get image file
			file, _, err := r.FormFile("image-upload")
			if err != nil {
				log.Error(err.Error())
				http.Error(w, "Image not found", http.StatusBadRequest)
				return
			}

			defer file.Close()

			adFolder := r.Form.Get("type")
			levelID := r.Form.Get("level-id")
			if adFolder == "" || levelID == "" {
				http.Error(w, "Missing type or levelID", http.StatusBadRequest)
				return
			}

			// Map type to number
			typeNum, err := database.IntFromAdType(database.AdType(adFolder))
			if err != nil {
				log.Error("Invalid ad type: %s", err.Error())
				http.Error(w, "Invalid ad type", http.StatusBadRequest)
				return
			}

			// Create target folder
			targetDir := filepath.Join("..", "ad_storage", adFolder)
			err = os.MkdirAll(targetDir, os.ModePerm)
			if err != nil {
				log.Error("Failed to get directory %s", err.Error())
				http.Error(w, "Failed to get directory", http.StatusInternalServerError)
				return
			}

			fileName := fmt.Sprintf("%s.webp", uid)
			dstPath := filepath.Join(targetDir, fileName)

			// Delete old image
			if _, err := os.Stat(dstPath); err == nil {
				log.Debug("Removing old ad image at %s", dstPath)
				os.Remove(dstPath)
			}

			dst, err := os.Create(dstPath)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}

			defer dst.Close()

			if _, err := io.Copy(dst, file); err != nil {
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Build a URL-ish path for retrieval; adjust base URL if you serve ads statically
			imageURL := strings.Join([]string{access.GetDomain(r), "cdn", adFolder, fileName}, "/")

			// Create DB row for the advertisement
			adID, err := database.CreateAdvertisement(uid, levelID, typeNum, imageURL)
			if err != nil {
				e := os.Remove(dstPath)
				if e != nil {
					log.Error("Failed to delete advertisement image: %s", e.Error())
				}

				log.Error("Failed to create advertisement row: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Info("Saved ad to %s, ad_id=%v, user_id=%s", dstPath, adID, uid)
			w.Write(fmt.Appendf(nil, `{"status":"ok","ad_id":%d,"image_url":"%s"}`, adID, imageURL))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
