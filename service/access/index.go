package access

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"service/database"
	"service/log"
)

func isInternal(ip string) bool {
	return strings.HasPrefix(ip, "127.") ||
		strings.HasPrefix(ip, "::1") ||
		strings.HasPrefix(ip, "192.168.") ||
		strings.HasPrefix(ip, "10.") ||
		strings.HasPrefix(ip, "172.")
}

// Check if the request was received internally (for testing with sensitive data)
func Restrict(ip string) (int, error) {
	log.Debug("Checking internal address %s", ip)

	if isInternal(ip) {
		log.Error("Address %s forbidden for use of internal API", ip)
		return http.StatusForbidden, fmt.Errorf("Forbidden")
	} else {
		log.Info("Address %s authorized for use of internal API", ip)
		return http.StatusOK, nil
	}
}

func GetDomain(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}

	return fmt.Sprintf("%s://%s", scheme, r.Host)
}

func FullURL(r *http.Request) string {
	base := GetDomain(r)
	return fmt.Sprintf("%s%s", base, r.RequestURI)
}

func init() {
	http.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			}

			if !u.IsAdmin {
				log.Error("User of ID %s is not admin", u.ID)
				http.Error(w, "User is not admin", http.StatusUnauthorized)
				return
			}

			users, err := database.GetAllUsers()
			if err != nil {
				log.Error("Failed to get all users: %s", err.Error())
				http.Error(w, "Failed to get all users", http.StatusInternalServerError)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(users); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/ban", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "POST")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodPost {
			header.Set("Content-Type", "application/json")

			// require login
			uid, err := GetSessionUserID(r)
			if err != nil || uid == "" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			} else {
				log.Info("Received user of ID %s", uid)
			}

			u, err := database.GetUser(uid)
			if err != nil {
				log.Error("Failed to get user: %s", err.Error())
				http.Error(w, "Failed to get user", http.StatusInternalServerError)
				return
			} else {
				log.Info("Fetched user %s (%s)", u.Username, u.ID)
			}

			if !u.IsAdmin {
				log.Error("User of ID %s is not admin", u.ID)
				http.Error(w, "User is not admin", http.StatusUnauthorized)
				return
			} else {
				log.Info("User %s (%s) is an admin", u.Username, u.ID)
			}

			query := r.URL.Query()
			idStr := query.Get("id")

			banned, err := database.BanUser(idStr)
			if err != nil {
				log.Error("Failed to ban user: %s", err.Error())
				http.Error(w, "Failed to ban user", http.StatusInternalServerError)
				return
			} else {
				log.Info("Banned user %s (%s)", banned.Username, banned.ID)
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(banned); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			} else {
				log.Info("Finished ban request to %s by admin %s (%s)", banned.Username, u.Username, u.ID)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
