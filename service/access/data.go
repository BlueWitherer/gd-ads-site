package access

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"service/log"

	_ "github.com/go-sql-driver/mysql"
)

var data *sql.DB

type AdEvent string

const ( // Table to save stats
	AdEventView  AdEvent = "views"  // For views
	AdEventClick AdEvent = "clicks" // For clicks
)

type StatBy string

const ( // Row to filter through
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

// private method to safely prepare the sql statement
func safePrepare(db *sql.DB, sql string) (*sql.Stmt, error) {
	if db != nil {
		log.Debug("Preparing connection for statement")
		return db.Prepare(sql)
	} else {
		return nil, fmt.Errorf("database connection non-existent")
	}
}

// Register a new client event for an ad
func NewStat(event AdEvent, ad int64, user int64) error {
	log.Debug("Registering new " + event)
	sql := fmt.Sprintf("INSERT INTO %s (ad_id, user_id, timestamp) VALUES (?, ?, ?)", event)

	stmt, err := safePrepare(data, sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ad, user, time.Now())
	return err
}

// inserts a new user or updates username if it already exists.
func UpsertUser(id string, username string) error {
	if id == "" {
		return fmt.Errorf("empty user id")
	}

	stmt, err := safePrepare(data, "INSERT INTO users (username, id) VALUES (?, ?) ON DUPLICATE KEY UPDATE username = VALUES (?), updated_at = CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, id)
	return err
}

// increments total_views or total_clicks for a user.
func IncrementUserStats(userID string, viewsDelta int, clicksDelta int) error {
	if userID == "" {
		return fmt.Errorf("empty user id")
	}

	stmt, err := safePrepare(data, "UPDATE users SET total_views = total_views + ?, total_clicks = total_clicks + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(viewsDelta, clicksDelta, userID)
	return err
}

// inserts an ad row
func CreateAdvertisement(userID, levelID string, adType int, imageURL string) (int64, error) {
	if userID == "" || levelID == "" || imageURL == "" {
		return 0, fmt.Errorf("missing ad fields")
	}

	stmt, err := safePrepare(data, "INSERT INTO advertisements (user_id, level_id, type, image_url) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(userID, levelID, adType, imageURL)
	if err != nil {
		return 0, err
	}

	id, _ := res.LastInsertId()
	return id, nil
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

func init() {
	var err error

	log.Info("Connecting to database with URI: " + os.Getenv("DB_URI"))
	data, err = sql.Open("mysql", os.Getenv("DB_URI"))
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = data.Ping()
	if err != nil {
		log.Error(err.Error())
		return
	} else if data == nil {
		log.Error("Database connection is nil")
		return
	}

	log.Print("MariaDB connection established.")
}
