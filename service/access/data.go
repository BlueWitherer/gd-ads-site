package access

import (
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	"service/log"

	_ "github.com/go-sql-driver/mysql"
)

var data *sql.DB

type AdEvent string

const (
	AdEventView  AdEvent = "view"  // For views
	AdEventClick AdEvent = "click" // For clicks
)

type StatBy string

const (
	StatByAd   StatBy = "ad_id"   // Filter stats by ad
	StatByUser StatBy = "user_id" // Filter stats by user
)

// Database row for advertisements listing
type AdRow struct {
	AdID     int64  `json:"ad_id"`
	UserID   string `json:"user_id"`
	LevelID  string `json:"level_id"`
	Type     int    `json:"type"`
	ImageURL string `json:"image_url"`
	Created  string `json:"created_at"`
}

// Register a new client event for an ad
func NewStat(event AdEvent, ad int64, user interface{}) error {
	log.Debug("Registering new " + event)
	stmt, err := data.Prepare("INSERT INTO ad_views (ad_id, user_id, timestamp) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	// Coerce numeric user to string transparently
	userID := user.(int64) // user is passed as int64
	_, err = stmt.Exec(ad, userID, time.Now())
	return err
}

func init() {
	var err error

	data, err = sql.Open("mysql", os.Getenv("DB_URI"))
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = data.Ping()
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Print("MariaDB connection established.")
}

// inserts a new user or updates username if it already exists.
func UpsertUser(id string, username string) error {
	if id == "" {
		return errors.New("empty user id")
	}
	// For MariaDB: use INSERT ... ON DUPLICATE KEY UPDATE
	q := `INSERT INTO users (id, username) VALUES (?, ?)
		  ON DUPLICATE KEY UPDATE username = VALUES(username), updated_at = CURRENT_TIMESTAMP`
	_, err := data.Exec(q, id, username)
	return err
}

// increments total_views or total_clicks for a user.
func IncrementUserStats(userID string, viewsDelta, clicksDelta int) error {
	if userID == "" {
		return errors.New("empty user id")
	}
	q := `UPDATE users SET total_views = total_views + ?, total_clicks = total_clicks + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`
	_, err := data.Exec(q, viewsDelta, clicksDelta, userID)
	return err
}

// inserts an ad row
func CreateAdvertisement(userID, levelID string, adType int, imageURL string) (int64, error) {
	if userID == "" || levelID == "" || imageURL == "" {
		return 0, errors.New("missing ad fields")
	}
	res, err := data.Exec(`INSERT INTO advertisements (user_id, level_id, type, image_url) VALUES (?, ?, ?, ?)`, userID, levelID, adType, imageURL)
	if err != nil {
		return 0, err
	}
	id, _ := res.LastInsertId()
	return id, nil
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
			err = errors.New("no user in session")
		}
		return "", err
	}
	return u.ID, nil
}

// fetches all ads for a given user
func ListAdvertisementsByUser(userID string) ([]AdRow, error) {
	rows, err := data.Query(`SELECT ad_id, user_id, level_id, type, image_url, created_at FROM advertisements WHERE user_id = ? ORDER BY ad_id DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []AdRow
	for rows.Next() {
		var r AdRow
		if err := rows.Scan(&r.AdID, &r.UserID, &r.LevelID, &r.Type, &r.ImageURL, &r.Created); err != nil {
			return nil, err
		}
		out = append(out, r)
	}
	return out, rows.Err()
}

// returns the owning user_id for an ad
func GetAdvertisementOwner(adID int64) (string, error) {
	var uid string
	err := data.QueryRow(`SELECT user_id FROM advertisements WHERE ad_id = ?`, adID).Scan(&uid)
	if err != nil {
		return "", err
	}
	return uid, nil
}
