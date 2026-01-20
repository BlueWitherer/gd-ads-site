package ads

import (
	"fmt"
	"net/http"
	"strconv"

	"service/access"
	"service/database"
	"service/discord"
	"service/log"
)

func init() {
	http.HandleFunc("/ads/boost", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := access.GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			query := r.URL.Query()

			idStr := query.Get("id")
			if idStr == "" {
				http.Error(w, "Missing ad ID parameter", http.StatusBadRequest)
				return
			}

			boostsStr := query.Get("boosts")
			if boostsStr == "" {
				http.Error(w, "Missing boosts parameter", http.StatusBadRequest)
				return
			}

			id, err := strconv.ParseInt(idStr, 10, 64)
			if err != nil {
				http.Error(w, "Invalid ad ID parameter", http.StatusBadRequest)
				return
			}

			boosts, err := strconv.ParseUint(boostsStr, 10, 32)
			if err != nil {
				http.Error(w, "Invalid boosts parameter", http.StatusBadRequest)
				return
			}

			user, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if int(user.BoostCount) < int(boosts) {
				http.Error(w, "Insufficient boosts", http.StatusBadRequest)
				return
			}

			ad, err := database.BoostAd(id, uint(boosts), user.ID)
			if err != nil {
				log.Error("Failed to boost advertisement: %s", err.Error())
				http.Error(w, "Failed to boost advertisement", http.StatusInternalServerError)
				return
			}

			err = discord.WebhookBoost(ad)
			if err != nil {
				log.Warn(err.Error())
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Successfully boosted ad")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
