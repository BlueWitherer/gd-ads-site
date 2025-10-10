package api

import (
	"net/http"
	"strconv"

	"service/access"
	"service/log"
)

func init() {
	http.HandleFunc("/api/view", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Registering view...")
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
			log.Error("Failed to get ad ID")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		user, err = strconv.ParseInt(userStr, 10, 64)
		if err != nil {
			log.Error("Failed to get user ID")
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = access.NewStat(access.AdEventView, ad, user)
		if err != nil {
			log.Error(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		} else {
			// Increment ad owner's total views
			if ownerID, ownerErr := access.GetAdvertisementOwner(ad); ownerErr == nil && ownerID != "" {
				if incErr := access.IncrementUserStats(ownerID, 1, 0); incErr != nil {
					log.Error("Failed to increment total views: " + incErr.Error())
				}
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("View registered!"))
		}
	})
}
