package api

import (
	"net/http"
	"strconv"

	"service/database"
	"service/log"
)

func init() {
	http.HandleFunc("/api/click", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Registering click...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		var ad int64   // Ad ID
		var user int64 // User ID

		var err error

		query := r.URL.Query()
		adStr := query.Get("ad_id")
		userStr := query.Get("user_id")

		ad, err = strconv.ParseInt(adStr, 10, 64)
		if err != nil {
			log.Error("Failed to get ad ID: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err = strconv.ParseInt(userStr, 10, 64)
		if err != nil {
			log.Error("Failed to get user ID: %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = database.NewStat(database.AdEventClick, ad, user)
		if err != nil {
			log.Error("Failed to create database click statistic: %s", err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			if ownerID, ownerErr := database.GetAdvertisementOwner(ad); ownerErr == nil && ownerID != "" {
				if incErr := database.IncrementUserStats(ownerID, 0, 1); incErr != nil {
					log.Error("Failed to increment total clicks: %s", incErr.Error())
				}
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Click registered!"))
		}
	})
}
