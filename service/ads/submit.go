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

		header.Set("Content-Type", "application/json")

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
		typeNum := 0
		switch adFolder {
		case string(database.AdTypeBanner):
			typeNum = 1

		case string(database.AdTypeSquare):
			typeNum = 2

		case string(database.AdTypeSkyscraper):
			typeNum = 3

		default:
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

		fileName := filepath.Base(handler.Filename)
		dstPath := filepath.Join(targetDir, fileName)
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
		adID, err := database.CreateAdvertisement(userID, levelID, typeNum, imageURL)
		if err != nil {
			log.Error("Failed to create advertisement row: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Info("Saved image to %s, ad_id=%v, user_id=%s", dstPath, adID, userID)
		w.Write(fmt.Appendf(nil, `{"status":"ok","ad_id":%d,"image_url":"%s"}`, adID, imageURL))
	})
}
