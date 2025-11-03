package access

import (
	"context"
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
	"service/utils"

	"github.com/google/uuid"
	"github.com/patrickmn/go-cache"
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

var sessionCache = cache.New(2*time.Hour, 10*time.Minute)

func generateSessionID() string {
	return uuid.New().String() // uuid v4
}

func isSecure(r *http.Request) bool {
	if r.TLS != nil || os.Getenv("ENV") == "production" {
		return true
	}

	return false
}

func SetSession(w http.ResponseWriter, user DiscordUser, secure bool) (string, error) {
	sessionId := generateSessionID()
	session := &http.Cookie{
		Name:     "session_id",
		Value:    sessionId,
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

	stmt, err := utils.PrepareStmt(utils.Db(), "INSERT INTO sessions (session_id, user_id, username, discriminator, avatar) VALUES (?, ?, ?, ?, ?) ON DUPLICATE KEY UPDATE user_id = VALUES(user_id), username = VALUES(username), discriminator = VALUES(discriminator), avatar = VALUES(avatar);")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(sessionId, user.ID, user.Username, user.Discriminator, user.Avatar)
	if err != nil {
		return "", err
	}

	log.Debug("Setting session cookie...")
	http.SetCookie(w, session)

	sessionCache.Set(sessionId, user, cache.DefaultExpiration)

	return sessionId, nil
}

func GetSessionFromId(id string) (*DiscordUser, error) {
	var user DiscordUser
	if val, found := sessionCache.Get(id); found {
		user = val.(DiscordUser)
	}

	stmt, err := utils.PrepareStmt(utils.Db(), "SELECT user_id, username, discriminator, avatar FROM sessions WHERE session_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(id).Scan(&user.ID, &user.Username, &user.Discriminator, &user.Avatar)
	if err != nil {
		return nil, err
	}

	updStmt, err := utils.PrepareStmt(utils.Db(), "UPDATE sessions SET last_seen = CURRENT_TIMESTAMP WHERE session_id = ?")
	if err != nil {
		return nil, err
	}
	defer updStmt.Close()

	_, err = updStmt.Exec(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

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

func getAvatarURL(userID, avatarHash string) string {
	if avatarHash == "" {
		var avId int
		if len(userID) > 0 {
			lastChar := userID[len(userID)-1]
			avId = int(lastChar) % 5
		} else {
			avId = 0
		}

		return fmt.Sprintf("https://cdn.discordapp.com/embed/avatars/%d.png", avId)
	}

	return fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.webp", userID, avatarHash)
}

func CleanupExpiredSessions() error {
	stmt, err := utils.PrepareStmt(utils.Db(), "DELETE FROM sessions WHERE last_seen < NOW() - INTERVAL 30 DAY")
	if err != nil {
		return err
	}
	defer stmt.Close()

	res, err := stmt.Exec()
	if err != nil {
		return err
	}

	rowsAffected, _ := res.RowsAffected()
	log.Info("Expired sessions cleaned: %d", rowsAffected)

	return nil
}

var sessionCancel context.CancelFunc

func StopSessionCleanup() {
	if sessionCancel != nil {
		sessionCancel()
	}
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

			if err := database.UpsertUser(user.ID, user.Username, getAvatarURL(user.ID, user.Avatar)); err != nil {
				log.Error("Failed to upsert user: %s", err.Error())
				http.Error(w, "Failed to upsert user", http.StatusInternalServerError)
				return
			}

			log.Debug("Setting session...")
			sessionId, err := SetSession(w, user, isSecure(r))
			if err != nil {
				log.Error("Failed to set the user's session: %s", err.Error())
				http.Error(w, "Failed to set the user's session", http.StatusInternalServerError)
				return
			}

			if jb, err := json.Marshal(user); err != nil {
				log.Debug("Creating session: id=%s (failed to marshal user)", sessionId)
			} else {
				log.Debug("Creating session: id=%s user=%s", sessionId, string(jb))
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

			header.Set("Content-Type", "application/json")
			if jb, err := json.Marshal(user); err == nil {
				log.Debug("/session returning user: %s", string(jb))
			} else {
				log.Debug("/session returning user: (failed to marshal)")
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(user); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/logout", func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err == nil {
			sessionCache.Delete(cookie.Value)
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

	http.HandleFunc("/account/me", func(w http.ResponseWriter, r *http.Request) {
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

			// if user is banned
			if u.Banned {
				http.Error(w, "User is banned", http.StatusForbidden)
				return
			}

			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(u); err != nil {
				log.Error("Failed to encode response: %s", err.Error())
				http.Error(w, "Failed to encode response", http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	ctx, cancel := context.WithCancel(context.Background())
	sessionCancel = cancel

	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Info("Session sweeper stopped.")
				return
			default:
				if err := CleanupExpiredSessions(); err != nil {
					log.Error("Failed to clean up sessions: %s", err.Error())
				}

				time.Sleep(3 * time.Hour)
			}
		}
	}()
}
