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
	Quantity       uint   `json:"quantity"`
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

func boostCode(env string) (string, error) {
	code := os.Getenv(env)
	if code == "" {
		return "", fmt.Errorf("direct link code env %s is not defined!", env)
	} else {
		return code, nil
	}
}

func getBoostReward(code string) uint {
	codeBoost, err := boostCode("KOFI_LINK_BOOST")
	if err != nil {
		log.Error(err.Error())
		return 0
	}

	codeBoostOverdrive, err := boostCode("KOFI_LINK_BOOST_OVERDRIVE")
	if err != nil {
		log.Error(err.Error())
		return 0
	}

	switch code {
	case codeBoost:
		return 5

	case codeBoostOverdrive:
		return 50

	default:
		return 0
	}
}

func init() {
	http.HandleFunc("/api/order", func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Ko-fi webhook called")

		header := w.Header()

		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

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

			switch body.Type {
			case KofiTypeShopOrder:
				log.Debug("Processing Ko-fi shop order for user of ID %s...", body.DiscordUserID)

				for _, item := range body.ShopItems {
					if b := getBoostReward(item.DirectLinkCode); b > 0 {
						if err := database.AddBoostsToUser(body.DiscordUserID, item.Quantity*b); err != nil {
							log.Error("Failed to add boosts: %s", err.Error())
							http.Error(w, "Failed to add boosts", http.StatusInternalServerError)
							return
						}

						log.Info("Added %d boosts to user of ID %s", b, body.DiscordUserID)
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

				err = database.AddBoostsToUser(user.ID, 3)
				if err != nil {
					log.Error("Failed to add boosts: %s", err.Error())
					http.Error(w, "Failed to add boosts", http.StatusInternalServerError)
					return
				}

				if body.IsSubscriptionPayment {
					log.Info("Verified %s with subscription!", user.Username)
				} else {
					log.Warn("Unverified %s due to subscription failure", user.Username)
				}

			default:
				log.Error("Invalid payment type")
				http.Error(w, "Invalid payment type", http.StatusBadRequest)
				return
			}

			w.WriteHeader(http.StatusOK)
			fmt.Fprint(w, "Ko-fi webhook received and processed")
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
