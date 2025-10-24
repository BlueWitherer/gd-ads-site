package api

import (
	"encoding/json"
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
			header.Set("Content-Type", "application/json")
			header.Set("Cache-Control", "no-store")

			var adFolder database.AdType

			query := r.URL.Query()
			adTypeStr := query.Get("type")

			typeNum, err := strconv.Atoi(adTypeStr)
			if err != nil {
				log.Error("Failed to get ad type ID: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			adFolder, err = database.AdTypeFromInt(typeNum)
			if err != nil {
				log.Error("Failed to get ad folder: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
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

			if len(ads) <= 0 {
				log.Info("No ads found for type %s", adFolder)
				http.Error(w, "No ads found", http.StatusNotFound)
				return
			}

			log.Debug("Getting random %s type ad...", adFolder)
			i := rand.Intn(len(ads))
			ad := ads[i]

			// Get view and click stats for this ad
			views, clicks, err := database.GetAdStats(ad.AdID)
			if err != nil {
				log.Error("Failed to get ad stats: %s", err.Error())
			} else {
				ad.ViewCount = views
				ad.ClickCount = clicks
			}

			log.Info("Returning ad as JSON: %s", ad.ImageURL)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ad); err != nil {
				log.Error("Failed to encode ad response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/ad/cdn", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting random ad image...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "image/webp")
			header.Set("Cache-Control", "no-store")

			query := r.URL.Query()
			adTypeStr := query.Get("type")

			typeNum, err := strconv.Atoi(adTypeStr)
			if err != nil {
				log.Error("Failed to get ad type ID: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			adFolder, err := database.AdTypeFromInt(typeNum)
			if err != nil {
				log.Error("Failed to get ad folder: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
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

			if len(ads) <= 0 {
				log.Info("No ads found for type %s", adFolder)
				http.Error(w, "No ads found", http.StatusNotFound)
				return
			}

			log.Debug("Getting random %s type ad...", adFolder)
			i := rand.Intn(len(ads))
			ad := ads[i]

			if ad.ImageURL == "" {
				log.Error("Ad has no image URL")
				http.Error(w, "Ad has no image", http.StatusInternalServerError)
				return
			}

			// Fetch image from CDN/source and stream it back
			resp, err := http.Get(ad.ImageURL)
			if err != nil {
				log.Error("Failed to fetch image: %s", err.Error())
				http.Error(w, "Failed to fetch image", http.StatusBadGateway)
				return
			}
			defer resp.Body.Close()

			if resp.StatusCode != http.StatusOK {
				log.Error("Image fetch returned status: %d", resp.StatusCode)
				http.Error(w, "Failed to fetch image", http.StatusBadGateway)
				return
			}

			log.Info("Returning image for ad: %s", ad.ImageURL)
			w.WriteHeader(http.StatusOK)
			if _, err := io.Copy(w, resp.Body); err != nil {
				log.Error("Failed to write image response: %s", err.Error())
				// can't do much if streaming fails
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/ad/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting ad by id...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")
			header.Set("Cache-Control", "no-store")

			query := r.URL.Query()
			idStr := query.Get("id")

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				log.Error("Failed to get ad ID: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			ad, err := database.GetAdvertisement(id)
			if err != nil {
				log.Error("Failed to get ad: %s", err.Error())
				http.Error(w, "Failed to get ad", http.StatusInternalServerError)
				return
			}

			// Get view and click stats for this ad
			views, clicks, err := database.GetAdStats(ad.AdID)
			if err != nil {
				log.Error("Failed to get ad stats: %s", err.Error())
			} else {
				ad.ViewCount = views
				ad.ClickCount = clicks
			}

			log.Info("Returning ad as JSON: %s", ad.ImageURL)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ad); err != nil {
				log.Error("Failed to encode ad response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		}
	})
}
