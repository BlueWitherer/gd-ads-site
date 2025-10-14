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
		if r.Method == http.MethodPost {
			// Require logged-in session
			userID, err := access.GetSessionUserID(r)
			if err != nil || userID == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Parse form with 10MB limit
			r.ParseMultipartForm(10 << 20)

			// Get image file
			file, handler, err := r.FormFile("image-upload")
			if err != nil {
				http.Error(w, "Image not found", http.StatusBadRequest)
				return
			}

			defer file.Close()

			adFolder := r.Form.Get("type")
			levelID := r.Form.Get("levelID")
			if adFolder == "" || levelID == "" {
				http.Error(w, "Missing type or levelID", http.StatusBadRequest)
				return
			}

			// Map type to number
			typeNum := 0
			switch adFolder {
			case "banner":
				typeNum = 1

			case "square":
				typeNum = 2

			case "skyscraper":
				typeNum = 3

			default:
				http.Error(w, "Invalid ad type", http.StatusBadRequest)
				return
			}

			// Create target folder
			targetDir := filepath.Join("..", "..", "ad_storage", adFolder)
			err = os.MkdirAll(targetDir, os.ModePerm)
			if err != nil {
				log.Error("Failed to get directory %s", err.Error())
				http.Error(w, "Failed to get directory", http.StatusInternalServerError)
				return
			}

			// Save file
			// Sanitize filename to prevent path traversal
			baseName := filepath.Base(handler.Filename)
			dstPath := filepath.Join(targetDir, baseName)
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
			imageURL := strings.Join([]string{"/ads", adFolder, baseName}, "/")

			// Create DB row for the advertisement
			adID, err := database.CreateAdvertisement(userID, levelID, typeNum, imageURL)
			if err != nil {
				log.Error("Failed to create advertisement row: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			log.Info("Saved image to %s, ad_id=%v, user_id=%s", dstPath, adID, userID)
			w.Header().Set("Content-Type", "application/json")
			w.Write([]byte(fmt.Sprintf(`{"status":"ok","ad_id":%d,"image_url":"%s"}`, adID, imageURL)))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
	})
}
