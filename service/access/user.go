package access

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"service/database"
	"service/log"

	"github.com/google/uuid"
)

type DiscordUser struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

var sessions = map[string]DiscordUser{} // session ID -> DiscordUser{}

// Get an ongoing user session if found
func GetSessionFromId(id string) (*DiscordUser, error) {
	user, ok := sessions[id]
	if !ok {
		return nil, fmt.Errorf("user not found")
	}

	return &user, nil
}

// extracts the logged-in user's ID from the session cookie via access session map
func GetSessionUserID(r *http.Request) (string, error) {
	c, err := r.Cookie("session_id")
	if err != nil {
		return "", err
	}

	u, err := GetSessionFromId(c.Value)
	if err != nil || u == nil {
		if err == nil {
			err = fmt.Errorf("no user in session")
		}

		return "", err
	}

	return u.ID, nil
}

func GetSession(r *http.Request) (*DiscordUser, error) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return nil, err
	}

	user, err := GetSessionFromId(cookie.Value)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func generateSessionID() string {
	return uuid.New().String() // uuid v4
}

func init() {
	log.Info("Starting authorization handlers...")

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			redirectURL := "https://discord.com/oauth2/authorize?client_id=" + os.Getenv("DISCORD_CLIENT_ID") +
				"&redirect_uri=" + os.Getenv("DISCORD_REDIRECT_URI") +
				"&response_type=code" +
				"&scope=identify"

			user, err := GetSessionUserID(r)
			if err != nil {
				log.Error(err.Error())
				http.Redirect(w, r, redirectURL, http.StatusFound)
			} else if user != "" {
				log.Info("Redirecting from login to dashboard")
				http.Redirect(w, r, "/dashboard", http.StatusFound)
			} else {
				log.Debug("panic time ig")
				http.Redirect(w, r, redirectURL, http.StatusFound)
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		header := w.Header()

		header.Set("Access-Control-Allow-Origin", "*")
		header.Set("Access-Control-Allow-Methods", "GET")
		header.Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == http.MethodGet {
			code := r.URL.Query().Get("code")
			if code == "" {
				http.Error(w, "Missing code", http.StatusBadRequest)
				return
			}

			log.Info("Received Discord auth code %s", code)

			data := url.Values{}
			data.Set("client_id", os.Getenv("DISCORD_CLIENT_ID"))
			data.Set("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
			data.Set("grant_type", "authorization_code")
			data.Set("code", code)
			data.Set("redirect_uri", os.Getenv("DISCORD_REDIRECT_URI"))

			encoded := data.Encode()

			req, _ := http.NewRequest(http.MethodPost, "https://discord.com/api/oauth2/token", strings.NewReader(encoded))
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			log.Debug("Sending data request to Discord...")

			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, "Token exchange failed", http.StatusInternalServerError)
				return
			}

			defer resp.Body.Close()

			tokenResp := Token{}

			tokenBody, _ := io.ReadAll(resp.Body)
			log.Debug("Token endpoint status: %s", resp.Status)

			if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
				log.Error("Discord returned non-JSON: %s", string(tokenBody))
				http.Error(w, "Discord returned unexpected response", http.StatusInternalServerError)
				return
			}

			if resp.Request != nil {
				log.Debug("Token endpoint final URL: %s", resp.Request.URL.String())
			}

			if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
				log.Error("Failed to decode token response: %s", err.Error())
				http.Error(w, "Token decode failed", http.StatusInternalServerError)
				return
			}

			if tokenResp.AccessToken == "" {
				log.Error("No access token returned from Discord")
				http.Error(w, "No access token", http.StatusInternalServerError)
				return
			}

			// Fetch user info from Discord
			req, _ = http.NewRequest(http.MethodGet, "https://discord.com/api/users/@me", nil)
			req.Header.Set("Authorization", tokenResp.TokenType+" "+tokenResp.AccessToken)
			req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

			log.Debug("Getting user info...")
			resp, err = client.Do(req)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
				return
			}

			defer resp.Body.Close()

			user := DiscordUser{}

			userBody, _ := io.ReadAll(resp.Body)
			log.Debug("User endpoint status: %s", resp.Status)
			log.Debug("User endpoint body: %s", string(userBody))
			if err := json.Unmarshal(userBody, &user); err != nil {
				log.Error("Failed to decode user info: %s", err.Error())
				http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
				return
			}

			if user.ID == "" {
				log.Error("Discord returned empty user id")
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			u, err := database.GetUser(user.ID)
			if err != nil {
				log.Error(err.Error())
			} else if u.Banned {
				log.Error("User %s is banned", u.Username)
				http.Error(w, "User is banned", http.StatusForbidden)
				return
			}

			if err := database.UpsertUser(user.ID, user.Username); err != nil {
				log.Error("Failed to upsert user: %s", err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			sessionID := generateSessionID()
			sessions[sessionID] = user

			// Set cookie attributes depending on whether the connection is secure.
			secure := false
			if r.TLS != nil || os.Getenv("ENV") == "production" {
				secure = true
			}

			session := &http.Cookie{
				Name:     "session_id",
				Value:    sessionID,
				Path:     "/",
				HttpOnly: true,
				Secure:   secure,
				Expires:  time.Now().Add(30 * 24 * time.Hour),
			}

			if secure {
				session.SameSite = http.SameSiteNoneMode
			} else {
				session.SameSite = http.SameSiteLaxMode
			}

			log.Debug("Setting session cookie...")
			http.SetCookie(w, session)

			// Debug log created session
			if jb, err := json.Marshal(user); err != nil {
				log.Debug("Creating session: id=%s (failed to marshal user)", sessionID)
			} else {
				log.Debug("Creating session: id=%s user=%s", sessionID, string(jb))
			}

			log.Info("Redirecting to dashboard")
			http.Redirect(w, r, "/dashboard", http.StatusFound)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			if c, err := r.Cookie("session_id"); err == nil {
				log.Debug("/session request cookie: %s", c.Value)
			} else {
				log.Debug("/session request no cookie: %s", err.Error())
			}

			user, err := GetSession(r)
			if err != nil {
				log.Error(err.Error())
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}

			header := w.Header()

			w.WriteHeader(http.StatusOK)
			header.Set("Content-Type", "application/json")
			if jb, err := json.Marshal(user); err == nil {
				log.Debug("/session returning user: %s", string(jb))
			} else {
				log.Debug("/session returning user: (failed to marshal)")
			}

			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			delete(sessions, cookie.Value)
			log.Info("User %s logged out", cookie.Value)
		}

		secure := false
		if r.TLS != nil || os.Getenv("ENV") == "production" {
			secure = true
		}

		clearCookie := &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // bye bye cookie
			HttpOnly: true,
			Secure:   secure,
			SameSite: http.SameSiteNoneMode,
		}

		http.SetCookie(w, clearCookie)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged out successfully"))
	})
}
