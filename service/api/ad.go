package api

import (
	"io"
	"math/rand"
	"net/http"
	"strconv"

	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/api/ad", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting random ad...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "image/webp")
			header.Set("Cache-Control", "no-store")

			var adFolder database.AdType

			query := r.URL.Query()
			adTypeStr := query.Get("type")

			typeNum, err := strconv.ParseInt(adTypeStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ad type: " + err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			// Map type to number
			switch typeNum {
			case 1:
				adFolder = database.AdTypeBanner

			case 2:
				adFolder = database.AdTypeSquare

			case 3:
				adFolder = database.AdTypeSkyscraper

			default:
				http.Error(w, "Invalid ad type", http.StatusBadRequest)
				return
			}

			rows, err := database.ListAllAdvertisements()
			if err != nil {
				log.Error("Failed to list ads: %s", err.Error())
				http.Error(w, "Failed to list ads", http.StatusInternalServerError)
				return
			}

			log.Debug("Filtering for %s type ads...", adFolder)
			ads, err := database.FilterAdsByType(rows, adFolder)
			if err != nil {
				log.Error("Failed to filter through ads: %s", err.Error())
				http.Error(w, "Failed to filter through ads", http.StatusInternalServerError)
				return
			}

			log.Debug("Getting random %s type ad...", adFolder)
			i := rand.Intn(len(ads))
			ad := ads[i]

			log.Info("Rendering image %s...", ad.ImageURL)
			resp, err := http.Get(ad.ImageURL)
			if err != nil || resp.StatusCode != http.StatusOK {
				http.Error(w, "Failed to fetch image", http.StatusInternalServerError)
				return
			}
			defer resp.Body.Close()

			_, err = io.Copy(w, resp.Body)
			if err != nil {
				http.Error(w, "Failed to stream image", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
