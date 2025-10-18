package api

import (
	"encoding/json"
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

			log.Info("Returning ad as JSON: %s", ad.ImageURL)
			respData, err := json.Marshal(ad)
			if err != nil {
				log.Error("Failed to marshal ad to JSON: %s", err.Error())
				http.Error(w, "Failed to encode ad", http.StatusInternalServerError)
				return
			}

			_, err = w.Write(respData)
			if err != nil {
				log.Error("Failed to write response: %s", err.Error())
				// connection write failed; nothing more to do
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
