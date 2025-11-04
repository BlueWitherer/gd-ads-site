package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"service/access"
	"service/database"
	"service/log"
	"service/utils"
)

func newStat(r *http.Request, adEvent utils.AdEvent) (int, error) {
	var body struct {
		AdID      int64  `json:"ad_id"`
		UserID    string `json:"user_id"`
		AccountID int    `json:"account_id"`
		AuthToken string `json:"authtoken"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to parse JSON body: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	log.Debug("Body decoded - AdID: %v, UserID: %s", body.AdID, body.UserID)

	user := utils.ArgonUser{Account: body.AccountID, Token: body.AuthToken}
	valid, err := access.ValidateArgonUser(user)
	if err != nil {
		log.Error("Failed to validate Argon user: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	if valid {
		if body.UserID == "" {
			log.Error("User ID is empty")
			return http.StatusBadRequest, err
		}

		err := database.NewStat(adEvent, body.AdID, body.UserID)
		if err != nil {
			log.Error("Failed to create database click statistic: %s", err.Error())
			return http.StatusInternalServerError, err
		} else {
			log.Info("%s passed for player %d", adEvent, body.AccountID)
		}
	} else {
		return http.StatusUnauthorized, fmt.Errorf("argon user invalid")
	}

	return http.StatusOK, nil
}

func init() {
	http.HandleFunc("/api/click", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Registering click...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			status, err := newStat(r, utils.AdEventClick)
			if err != nil {
				http.Error(w, "Failed to register click statistic", status)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Click registered!"))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/api/view", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Registering view...")
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			status, err := newStat(r, utils.AdEventView)
			if err != nil {
				http.Error(w, "Failed to register view statistic", status)
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("View registered!"))
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
