package access

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"service/database"
	"service/log"
)

type ArgonValidation struct {
	Valid bool   `json:"valid"` // If user is valid
	Cause string `json:"cause"` // Cause for invalidation if any
}

type ArgonUser struct {
	Account int    `json:"account_id"` // Player account ID
	Token   string `json:"authtoken"`  // Authorization token
}

var ArgonCache []ArgonUser

func UpsertArgonUser(user ArgonUser) {
	for i, u := range ArgonCache {
		if u.Account == user.Account {
			ArgonCache[i].Token = user.Token
			return
		}
	}

	ArgonCache = append(ArgonCache, user)
}

func ValidateArgonUser(user ArgonUser) (bool, error) {
	for _, u := range ArgonCache {
		if u.Token == user.Token {
			log.Info("Argon cache hit for account of ID %v", user.Account)
			return true, nil
		}
	}

	u, err := url.Parse("https://argon.globed.dev/v1/validation/check")
	if err != nil {
		return false, err
	} else {
		log.Debug("Argon URL parsed for account of ID %v", user.Account)
	}

	q := u.Query()
	q.Set("account_id", fmt.Sprintf("%v", user.Account))
	q.Set("authtoken", user.Token)
	u.RawQuery = q.Encode()

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return false, err
	} else {
		log.Debug("Argon request object constructed for account of ID %v", user.Account)
	}

	req.Header.Set("User-Agent", "PlayerAdvertisements/1.0")

	log.Info("Sending request to Argon server: %s", u.String())
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	} else {
		log.Debug("Argon status received for account of ID %v", user.Account)
	}
	defer resp.Body.Close()

	var valid ArgonValidation
	if err := json.NewDecoder(resp.Body).Decode(&valid); err != nil {
		return false, err
	} else {
		log.Debug("Argon status of account of ID %v retrieved", user.Account)
	}

	if valid.Valid {
		log.Info("Argon status of account of ID %v is valid", user.Account)
		UpsertArgonUser(user)
		return true, nil
	}

	log.Error("Argon status of account of ID %v is invalid for %s", user.Account, valid.Cause)
	return false, fmt.Errorf("cause: %s", valid.Cause)
}

func DeleteArgonUser(accountID int) {
	for i, u := range ArgonCache {
		if u.Account == accountID {
			ArgonCache = append(ArgonCache[:i], ArgonCache[i+1:]...)
			return
		}
	}
}

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

	http.HandleFunc("/unban", func(w http.ResponseWriter, r *http.Request) {
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

			unbanned, err := database.UnbanUser(idStr)
			if err != nil {
				log.Error("Failed to unban user: %s", err.Error())
				http.Error(w, "Failed to unban user", http.StatusInternalServerError)
				return
			} else {
				log.Info("Unbanned user %s (%s)", unbanned.Username, unbanned.ID)
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(unbanned); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			} else {
				log.Info("Finished unban request to %s by admin %s (%s)", unbanned.Username, u.Username, u.ID)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			// require login and admin status
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

			// Extract user ID or username from URL path
			searchQuery := strings.TrimPrefix(r.URL.Path, "/users/")
			if searchQuery == "" {
				http.Error(w, "Missing user ID or username", http.StatusBadRequest)
				return
			}

			log.Debug("Searching for user: %s", searchQuery)

			// Try to get user by ID first
			targetUser, err := database.GetUser(searchQuery)
			if err != nil {
				// If not found by ID, try to find by username
				allUsers, err := database.GetAllUsers()
				if err != nil {
					log.Error("Failed to get all users: %s", err.Error())
					http.Error(w, "User not found", http.StatusNotFound)
					return
				}

				found := false
				for _, user := range allUsers {
					if user.Username == searchQuery {
						targetUser = user
						found = true
						break
					}
				}

				if !found {
					log.Error("User not found by ID or username: %s", searchQuery)
					http.Error(w, "User not found", http.StatusNotFound)
					return
				}
			}

			// Get all ads and filter by user
			allAds, err := database.ListAllAdvertisements()
			if err != nil {
				log.Error("Failed to list ads: %s", err.Error())
				http.Error(w, "Failed to list ads", http.StatusInternalServerError)
				return
			}

			userAds, err := database.FilterAdsByUser(allAds, targetUser.ID)
			if err != nil {
				log.Error("Failed to filter ads: %s", err.Error())
				http.Error(w, "Failed to filter ads", http.StatusInternalServerError)
				return
			}

			// Populate view/click counts for each ad
			for i := range userAds {
				views, clicks, err := database.GetAdStats(userAds[i].AdID)
				if err != nil {
					log.Debug("Failed to get ad stats: %s", err.Error())
				} else {
					userAds[i].ViewCount = views
					userAds[i].ClickCount = clicks
				}
			}

			response := map[string]interface{}{
				"user": targetUser,
				"ads":  userAds,
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}

			log.Info("Admin %s fetched user %s info", u.Username, targetUser.ID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/users/fetch", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type, User-Agent")

		if r.Method == http.MethodGet {
			header.Set("Content-Type", "application/json")

			// Check User-Agent header
			userAgent := r.Header.Get("User-Agent")
			if userAgent != "PlayerAdvertisements/1.0" {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			// Get user ID from query parameter
			query := r.URL.Query()
			searchQuery := query.Get("id")
			if searchQuery == "" {
				http.Error(w, "Missing user ID parameter", http.StatusBadRequest)
				return
			}

			log.Debug("Fetching user: %s", searchQuery)

			// Try to get user by ID first
			targetUser, err := database.GetUser(searchQuery)
			if err != nil {
				// If not found by ID, try to find by username
				allUsers, err := database.GetAllUsers()
				if err != nil {
					log.Error("Failed to get all users: %s", err.Error())
					http.Error(w, "User not found", http.StatusNotFound)
					return
				}

				found := false
				for _, user := range allUsers {
					if user.Username == searchQuery {
						targetUser = user
						found = true
						break
					}
				}

				if !found {
					log.Error("User not found by ID or username: %s", searchQuery)
					http.Error(w, "User not found", http.StatusNotFound)
					return
				}
			}

			// Get all ads and filter by user
			allAds, err := database.ListAllAdvertisements()
			if err != nil {
				log.Error("Failed to list ads: %s", err.Error())
				http.Error(w, "Failed to list ads", http.StatusInternalServerError)
				return
			}

			userAds, err := database.FilterAdsByUser(allAds, targetUser.ID)
			if err != nil {
				log.Error("Failed to filter ads: %s", err.Error())
				http.Error(w, "Failed to filter ads", http.StatusInternalServerError)
				return
			}

			// Populate view/click counts for each ad
			for i := range userAds {
				views, clicks, err := database.GetAdStats(userAds[i].AdID)
				if err != nil {
					log.Debug("Failed to get ad stats: %s", err.Error())
				} else {
					userAds[i].ViewCount = views
					userAds[i].ClickCount = clicks
				}
			}

			response := map[string]interface{}{
				"user": targetUser,
				"ads":  userAds,
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}

			log.Info("Fetched user %s info via /users/fetch", targetUser.ID)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
