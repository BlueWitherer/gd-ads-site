package ads

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

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

			user, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get ad owner: %s", err.Error())
				http.Error(w, "Failed to get ad owner", http.StatusInternalServerError)
				return
			}

			if user.Banned {
				log.Error("User %s is banned", user.Username)
				http.Error(w, "User is banned", http.StatusForbidden)
				return
			}

			// Check if user has reached the maximum number of active ads (10)
			activeAdCount, err := database.CountActiveAdvertisementsByUser(uid)
			if err != nil {
				log.Error("Failed to count active advertisements: %s", err.Error())
				http.Error(w, "Failed to check advertisement limit", http.StatusInternalServerError)
				return
			}

			if activeAdCount >= 10 {
				log.Warn("User %s attempted to submit ad but has reached maximum (10) active advertisements", user.Username)
				http.Error(w, "You have reached the maximum number of active advertisements (10)", http.StatusBadRequest)
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

			dst, err := os.Create(dstPath)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, "Failed to save image", http.StatusInternalServerError)
				return
			}

			if _, err := io.Copy(dst, file); err != nil {
				dst.Close()
				log.Error(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Close the file before renaming
			dst.Close()

			// Create DB row for the advertisement first (to get the ad ID)
			// Use a placeholder URL initially
			placeholderURL := fmt.Sprintf("%s/cdn/%s/placeholder?v=%d", access.GetDomain(r), adFolder, time.Now().Unix())
			adID, err := database.CreateAdvertisement(uid, levelID, typeNum, placeholderURL)
			if err != nil {
				e := os.Remove(dstPath)
				if e != nil {
					log.Error("Failed to delete advertisement image: %s", e.Error())
				}

				log.Error("Failed to create advertisement row: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			// Now rename the file to include the ad ID
			newFileName := fmt.Sprintf("%s-%d.webp", uid, adID)
			newDstPath := filepath.Join(targetDir, newFileName)
			err = os.Rename(dstPath, newDstPath)
			if err != nil {
				_, delErr := database.DeleteAdvertisement(adID)
				if delErr != nil {
					log.Error("Failed to delete advertisement row: %s", delErr.Error())
				}
				e := os.Remove(dstPath)
				if e != nil {
					log.Error("Failed to delete advertisement image: %s", e.Error())
				}

				log.Error("Failed to rename advertisement image: %s", err.Error())
				http.Error(w, "Failed to rename advertisement image", http.StatusInternalServerError)
				return
			}

			// Update the image URL with the correct filename
			imageURL := fmt.Sprintf("%s/cdn/%s/%s?v=%d", access.GetDomain(r), adFolder, newFileName, time.Now().Unix())
			err = database.UpdateAdvertisementImageURL(adID, imageURL)
			if err != nil {
				_, delErr := database.DeleteAdvertisement(adID)
				if delErr != nil {
					log.Error("Failed to delete advertisement row: %s", delErr.Error())
				}
				e := os.Remove(newDstPath)
				if e != nil {
					log.Error("Failed to delete advertisement image: %s", e.Error())
				}

				log.Error("Failed to update advertisement image URL: %s", err.Error())
				http.Error(w, "Failed to update advertisement image URL", http.StatusInternalServerError)
				return
			}

			if user.IsAdmin {
				newAd, err := database.ApproveAd(adID)
				if err != nil {
					log.Error("Failed to auto-approve new ad by admin: %s", err.Error())
				} else {
					log.Info("Auto-approved ad %s (%v) by admin %s (%s)", newAd.ImageURL, newAd.AdID, user.Username, user.ID)
				}
			}

			log.Info("Saved ad to %s, ad_id=%v, user_id=%s", newDstPath, adID, uid)
			w.Write(fmt.Appendf(nil, `{"status":"ok","ad_id":%d,"image_url":"%s"}`, adID, imageURL))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
