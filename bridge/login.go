package main

import (
	"bridge/log"
	"encoding/json"
	"net/http"
	"net/url"
	"os"
	"strings"

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

func getUserFromSession(r *http.Request) *User {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Error(err.Error())
		return nil
	}

	user, ok := sessions[cookie.Value]
	if !ok {
		log.Error("User " + cookie.Value + " not found")
		return nil
	}

	log.Info("User " + cookie.Value + " found and authorized")
	return &user
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

		req, _ := http.NewRequest("POST", "https://discord.com/api/oauth2/token", strings.NewReader(data.Encode()))
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

		log.Debug("Decoding token...")
		json.NewDecoder(resp.Body).Decode(&tokenResp)

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

		log.Debug("Decoding user info...")
		if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
			http.Error(w, "Failed to decode user info", http.StatusInternalServerError)
			return
		}

		sessionID := generateSessionID()
		sessions[sessionID] = user

		log.Debug("Setting session cookie...")
		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true,
			Secure:   true, // if using HTTPS
		})

		log.Info("Redirecting to dashboard")
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	})

	http.HandleFunc("/session", func(w http.ResponseWriter, r *http.Request) {
		user := getUserFromSession(r)

		if user == nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(user)
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			delete(sessions, cookie.Value)
			log.Info("User " + cookie.Value + " logged out")
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    "",
			Path:     "/",
			MaxAge:   -1, // bye bye cookie
			HttpOnly: true,
			Secure:   true,
		})

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Logged out successfully"))
	})
}
