package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"service/access"
	"service/database"
	"service/log"
	"service/utils"
)

func newStat(r *http.Request, query url.Values, adEvent utils.AdEvent) (int, error) {
	accountIDStr := query.Get("account_id")
	authToken := query.Get("authtoken")

	var body struct {
		AdID   int64  `json:"ad_id"`
		UserID string `json:"user_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		log.Error("Failed to parse JSON body: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	log.Debug("Body decoded - AdID: %v, UserID: %s", body.AdID, body.UserID)

	if accountIDStr == "" || authToken == "" {
		log.Error("Missing query parameters - account_id: %s, authtoken: %s", accountIDStr, authToken)
		return http.StatusBadRequest, fmt.Errorf("missing query parameters")
	}

	accountId, err := strconv.Atoi(accountIDStr)
	if err != nil {
		log.Error("Failed to parse account ID: %s", err.Error())
		return http.StatusInternalServerError, err
	}

	user := utils.ArgonUser{Account: accountId, Token: authToken}
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

		err := database.NewStatWithUserID(adEvent, body.AdID, body.UserID)
		if err != nil {
			log.Error("Failed to create database click statistic: %s", err.Error())
			return http.StatusInternalServerError, err
		} else {
			log.Info("click passed: %s", accountIDStr)
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
			query := r.URL.Query()

			status, err := newStat(r, query, utils.AdEventClick)
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
			query := r.URL.Query()

			status, err := newStat(r, query, utils.AdEventView)
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
