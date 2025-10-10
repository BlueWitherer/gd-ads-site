package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"bridge/log"

	"github.com/google/uuid"
)

type User struct {
	ID            string `json:"id"`
	Username      string `json:"username"`
	Discriminator string `json:"discriminator"`
	Avatar        string `json:"avatar"`
}

type Token struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
}

var sessions = map[string]User{} // sessionID -> User

func generateSessionID() string {
	return uuid.New().String() // or any secure random generator
}

func GetSessionFromId(id string) (*User, error) {
	user, ok := sessions[id]
	if !ok {
		log.Error("User " + id + " not found")
		return nil, fmt.Errorf("User not found")
	}

	log.Info("User " + id + " found and authorized")
	return &user, nil
}

func getUserFromSession(r *http.Request) *User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	user, err := GetSessionFromId(cookie.Value)
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	return user
}

func init() {
	log.Info("Starting authorization handlers...")

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		redirectURL := "https://discord.com/oauth2/authorize?client_id=" + os.Getenv("DISCORD_CLIENT_ID") +
			"&redirect_uri=" + os.Getenv("DISCORD_REDIRECT_URI") +
			"&response_type=code" +
			"&scope=identify"
		http.Redirect(w, r, redirectURL, http.StatusFound)
	})

	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		code := r.URL.Query().Get("code")
		if code == "" {
			http.Error(w, "Missing code", http.StatusBadRequest)
			return
		}

		log.Info("Received Discord auth code " + code)

		data := url.Values{}
		data.Set("client_id", os.Getenv("DISCORD_CLIENT_ID"))
		data.Set("client_secret", os.Getenv("DISCORD_CLIENT_SECRET"))
		data.Set("grant_type", "authorization_code")
		data.Set("code", code)
		data.Set("redirect_uri", os.Getenv("DISCORD_REDIRECT_URI"))

		encoded := data.Encode()
		// redact client_secret in logs
		secret := os.Getenv("DISCORD_CLIENT_SECRET")
		redacted := strings.ReplaceAll(encoded, secret, "<redacted>")

		log.Debug("Token request body (redacted): " + redacted)

		req, _ := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(encoded))
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

		log.Debug("Sending data request to Discord...")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Token exchange failed", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		tokenResp := Token{}

		// Read token response body so we can log on error and still attempt decode
		tokenBody, _ := io.ReadAll(resp.Body)
		log.Debug("Token endpoint status: " + resp.Status)
		log.Debug("Token endpoint body: " + string(tokenBody))
		if resp.Request != nil {
			log.Debug("Token endpoint final URL: " + resp.Request.URL.String())
		}
		if err := json.Unmarshal(tokenBody, &tokenResp); err != nil {
			log.Error("Failed to decode token response: " + err.Error())
			http.Error(w, "Token decode failed", http.StatusInternalServerError)
			return
		}

		if tokenResp.AccessToken == "" {
			log.Error("No access token returned from Discord")
			http.Error(w, "No access token", http.StatusInternalServerError)
			return
		}

		// Fetch user info from Discord
		req, _ = http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
		req.Header.Set("Authorization", tokenResp.TokenType+" "+tokenResp.AccessToken)

		log.Debug("Getting user info...")
		resp, err = client.Do(req)
		if err != nil {
			http.Error(w, "Failed to fetch user info", http.StatusInternalServerError)
			return
		}

		defer resp.Body.Close()

		user := User{}

		// Read user response body so we can log and validate
		userBody, _ := io.ReadAll(resp.Body)
		log.Debug("User endpoint status: " + resp.Status)
		log.Debug("User endpoint body: " + string(userBody))
		if err := json.Unmarshal(userBody, &user); err != nil {
			log.Error("Failed to decode user info: " + err.Error())
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}

		if user.ID == "" {
			log.Error("Discord returned empty user id")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		sessionID := generateSessionID()
		sessions[sessionID] = user

		// Log the created session and user for debugging
		if jb, err := json.Marshal(user); err == nil {
			log.Debug("Creating session: id=" + sessionID + " user=" + string(jb))
		} else {
			log.Debug("Creating session: id=" + sessionID + " (failed to marshal user)")
		}

		secure := false
		if r.TLS != nil || os.Getenv("ENV") == "production" {
			secure = true
		}

		cookie := &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   secure,
		}

		if secure {
			cookie.SameSite = http.SameSiteNoneMode
		} else {
			cookie.SameSite = http.SameSiteLaxMode
		}

		log.Debug("Setting session cookie: " + cookie.String())
		http.SetCookie(w, cookie)

		log.Info("Redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		// Log incoming cookie on session request
		if c, err := r.Cookie("session_id"); err == nil {
			log.Debug("/session request cookie: " + c.Value)
		} else {
			log.Debug("/session request no cookie: " + err.Error())
		}

		user := getUserFromSession(r)

		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		if user != nil {
			if jb, err := json.Marshal(user); err == nil {
				log.Debug("/session returning user: " + string(jb))
			} else {
				log.Debug("/session returning user: (failed to marshal)")
			}
		} else {
			log.Debug("/session returning nil user")
		}

		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			delete(sessions, cookie.Value)
			log.Info("User " + cookie.Value + " logged out")
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
		}
		if secure {
			clearCookie.SameSite = http.SameSiteNoneMode
		} else {
			clearCookie.SameSite = http.SameSiteLaxMode
		}

		http.SetCookie(w, clearCookie)

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged out successfully"))
	})
}
