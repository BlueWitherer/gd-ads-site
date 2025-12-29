package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"service/database"
	"service/log"
)

type KofiType string

const (
	KofiTypeDonation     KofiType = "Donation"
	KofiTypeSubscription KofiType = "Subscription"
	KofiTypeCommission   KofiType = "Commission"
	KofiTypeShopOrder    KofiType = "Shop Order"
)

type KofiShopItem struct {
	DirectLinkCode string `json:"direct_link_code"`
	ItemName       string `json:"item_name"`
	Quantity       int    `json:"quantity"`
}

type Kofi struct {
	VerificationToken     string         `json:"verification_token"`
	Amount                string         `json:"amount"`
	Timestamp             time.Time      `json:"timestamp"`
	Type                  KofiType       `json:"type"`
	FromName              string         `json:"from_name"`
	Message               string         `json:"message"`
	IsSubscriptionPayment bool           `json:"is_subscription_payment"`
	IsPublic              bool           `json:"is_public"`
	ShopItems             []KofiShopItem `json:"shop_items"`
	DiscordUserID         string         `json:"discord_userid"`
}

func boostLinkCode() (string, error) {
	code := os.Getenv("KOFI_LINK_BOOST")
	if code == "" {
		return code, fmt.Errorf("direct link code for boost is not defined!")
	} else {
		return code, nil
	}
}

func init() {
	http.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Ko-fi webhook called")

		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			if err := r.ParseForm(); err != nil {
				log.Error("Failed to parse form data: %s", err.Error())
				http.Error(w, "Failed to parse form data", http.StatusBadRequest)
				return
			}

			log.Debug("Ko-fi values: %+v", r.Form)

			data := r.FormValue("data")
			if data == "" {
				log.Error("Missing form data")
				http.Error(w, "Missing form data", http.StatusBadRequest)
				return
			}

			var body Kofi
			if err := json.Unmarshal([]byte(data), &body); err != nil {
				log.Error("Failed to unmarshal JSON: %s", err.Error())
				http.Error(w, "Failed to unmarshal JSON", http.StatusBadRequest)
				return
			}

			if os.Getenv("KOFI_VERIFICATION_TOKEN") != body.VerificationToken {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			boostCode, err := boostLinkCode()
			if err != nil {
				log.Error("Failed to get Boost product code: %s", err.Error())
				http.Error(w, "Failed to get Boost product code", http.StatusBadRequest)
				return
			}

			switch body.Type {
			case KofiTypeShopOrder:
				log.Debug("Processing Ko-fi shop order for user of ID %s...", body.DiscordUserID)

				for _, item := range body.ShopItems {
					if item.DirectLinkCode == boostCode {
						if err := database.AddBoostsToUser(body.DiscordUserID, uint(item.Quantity)*5); err != nil {
							log.Error("Failed to add boosts: %s", err.Error())
							http.Error(w, "Failed to add boosts", http.StatusInternalServerError)
							return
						}

						log.Info("Added boosts to user of ID %s", body.DiscordUserID)
					}
				}

			case KofiTypeSubscription:
				log.Debug("Processing Ko-fi subscription for user of ID %s...", body.DiscordUserID)

				user, err := database.VerifyUser(body.DiscordUserID, body.IsSubscriptionPayment)
				if err != nil {
					log.Error("Failed to verify user through subscription: %s", err.Error())
					http.Error(w, "Failed to verify user through subscription", http.StatusInternalServerError)
					return
				}

				if body.IsSubscriptionPayment {
					log.Info("Verified %s with subscription!", user.Username)
				} else {
					log.Warn("Unverified %s due to subscription failure", user.Username)
				}

			default:
				log.Error("Invalid payment type")
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Ko-fi webhook received")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
