package api

import (
	"encoding/json"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"service/database"
	"service/log"
	"service/utils"

	"github.com/patrickmn/go-cache"
)

var globalStats = cache.New(10*time.Minute, 15*time.Minute)

func init() {
	http.HandleFunc("/api/ad", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Getting random ad...")
		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")
			header.Set("Cache-Control", "no-store")

			var adFolder utils.AdType

			query := r.URL.Query()
			adTypeStr := query.Get("type")

			typeNum, err := strconv.Atoi(adTypeStr)
			if err != nil {
				log.Error("Failed to get ad type ID: %s", err.Error())
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			adFolder, err = utils.AdTypeFromInt(typeNum)
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

			safeAds, err := database.FilterAdsFromBannedUsers(rows)
			if err != nil {
				log.Error("Failed to filter safe ads: %s", err.Error())
				http.Error(w, "Failed to filter safe ads", http.StatusInternalServerError)
				return
			}

			liveAds, err := database.FilterAdsByPending(safeAds, false)
			if err != nil {
				log.Error("Failed to filter pending ads: %s", err.Error())
				http.Error(w, "Failed to filter pending ads", http.StatusInternalServerError)
				return
			}

			log.Debug("Filtering for %s type ads...", adFolder)
			ads, err := database.FilterAdsByType(liveAds, adFolder)
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
			totalWeight := 0.0
			weights := make([]float64, len(ads))
			for idx, a := range ads {
				w := 1.0
				globalClicks := uint64(1)

				if val, found := globalStats.Get("global_clicks"); found {
					globalClicks = val.(uint64)
				} else {
					stats, err := database.GetGlobalStats()
					if err != nil {
						log.Error("Failed to get global ad stats: %s", err.Error())
					} else {
						globalClicks = uint64(stats.TotalClicks)
						globalStats.Set("global_clicks", globalClicks, cache.DefaultExpiration)
					}
				}

				if a.BoostCount > 0 {
					w += float64(a.BoostCount)
				}

				u, err := database.GetUser(a.UserID)
				if err != nil {
					log.Error("Failed to get ad owner for boosting: %s", err.Error())
				} else {
					if u.IsAdmin {
						w += 1
					}

					if u.IsStaff {
						w += 2
					}

					if u.Verified {
						w += 3
					}
				}

				if time.Since(a.Created).Hours() < 36 {
					denom := 0.05 * float64(globalClicks)
					if denom < 1 {
						denom = 1
					}

					w += 3 * math.Exp(-float64(a.Clicks)/denom)
				}

				w += float64(a.Clicks) * 0.125

				weights[idx] = w
				totalWeight += w
			}

			var chosenIdx int
			if totalWeight <= 0 {
				chosenIdx = rand.Intn(len(ads))
			} else {
				rn := rand.Float64() * totalWeight
				cn := 0.0
				for idx, w := range weights {
					cn += w
					if rn < cn {
						chosenIdx = idx
						break
					}
				}
			}
			ad := ads[chosenIdx]

			// Get view and click stats for this ad
			views, clicks, err := database.GetAdStats(ad.AdID)
			if err != nil {
				log.Error("Failed to get ad stats: %s", err.Error())
			} else {
				ad.Views = uint64(views)
				ad.Clicks = uint64(clicks)
			}

			log.Debug("Returning ad as JSON: %s", ad.ImageURL)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ad); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
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

			user, err := database.GetUser(ad.UserID)
			if err != nil {
				log.Error("Failed to get ad owner: %s", err.Error())
				http.Error(w, "Failed to get ad owner", http.StatusInternalServerError)
				return
			}

			if user.Banned {
				log.Error("Owner %s of advertisement of ID %v is banned", user.Username, ad.AdID)
				http.Error(w, "Advertisement owner is banned", http.StatusForbidden)
				return
			}

			// Get view and click stats for this ad
			views, clicks, err := database.GetAdStats(ad.AdID)
			if err != nil {
				log.Error("Failed to get ad stats: %s", err.Error())
			} else {
				ad.Views = uint64(views)
				ad.Clicks = uint64(clicks)
			}

			log.Info("Returning ad as JSON: %s", ad.ImageURL)
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(ad); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
