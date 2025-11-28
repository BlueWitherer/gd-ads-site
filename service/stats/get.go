package stats

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"service/access"
	"service/database"
	"service/log"

	"github.com/patrickmn/go-cache"
)

var downloads = cache.New(5*time.Minute, 10*time.Minute)

func init() {
	http.HandleFunc("/stats/get", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting advertisement stats for user...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			// Require logged-in user
			uid, err := access.GetSessionUserID(r)
			if err != nil {
				log.Error("Failed to get session ID: %s", err.Error())
				http.Error(w, "Failed to get session ID", http.StatusInternalServerError)
				return
			} else if uid == "" {
				log.Error("Unauthorized access to /stats/get")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			header.Set("Content-Type", "application/json")

			stats, err := database.GetUserTotals(uid)
			if err != nil {
				log.Error("Failed to fetch user totals: %s", err.Error())
				http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
				return
			}

			log.Info("Retrieved stats for user: %s", uid)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(stats); err != nil {
				log.Error("Failed to encode stats response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/stats/global", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting global advertisement statistics...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			stats, err := database.GetGlobalStats()
			if err != nil {
				log.Error("Failed to fetch global stats: %s", err.Error())
				http.Error(w, "Failed to fetch stats", http.StatusInternalServerError)
				return
			}

			log.Debug("Retrieved global stats - Views: %d, Clicks: %d, Ads: %d", stats.TotalViews, stats.TotalClicks, stats.AdCount)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(stats); err != nil {
				log.Error("Failed to encode global stats response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// sends get req to https://api.geode-sdk.org/v1/mods/arcticwoof.player_advertisements
	http.HandleFunc("/stats/downloads", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting download count for Player Advertisements on Geode...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			var count uint64 = 0

			if val, found := downloads.Get("count"); found {
				c := val.(uint64)

				log.Debug("Returning cached download count of %d", c)
				count = c
			} else {
				dl, err := database.GetModDownloads()
				if err != nil {
					log.Error("Failed to fetch mod download count: %s", err.Error())
					http.Error(w, "Failed to fetch mod download count", http.StatusInternalServerError)
					return
				}

				log.Info("New mod download count of %d", dl)

				downloads.Set("count", dl, cache.DefaultExpiration)
				count = dl
			}

			fmt.Fprintf(w, "%d", count)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
