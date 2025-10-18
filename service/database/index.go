package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"service/log"

	_ "github.com/go-sql-driver/mysql"
)

// Concurrent database connection
var data *sql.DB

type AdType string // Dimensions of the ad image

const ( // Table to save stats
	AdTypeBanner     AdType = "banner"     // Horizontal ads
	AdTypeSquare     AdType = "square"     // Square ads
	AdTypeSkyscraper AdType = "skyscraper" // Vertical ads
)

type AdEvent string // Table to save stats to

const (
	AdEventView  AdEvent = "views"  // For views
	AdEventClick AdEvent = "clicks" // For clicks
)

type StatBy string // Row to filter through

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

// private method to safely prepare the sql statement
func prepareStmt(db *sql.DB, sql string) (*sql.Stmt, error) {
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

	stmt, err := prepareStmt(data, sql)
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

	stmt, err := prepareStmt(data, "INSERT INTO users (username, id) VALUES (?, ?) ON DUPLICATE KEY UPDATE username = VALUES (username), updated_at = CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(username, id)
	return err
}

// increments total_views or total_clicks for a user.
func IncrementUserStats(userId string, viewsDelta int, clicksDelta int) error {
	if userId == "" {
		return fmt.Errorf("empty user id")
	}

	stmt, err := prepareStmt(data, "UPDATE users SET total_views = total_views + ?, total_clicks = total_clicks + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(viewsDelta, clicksDelta, userId)
	return err
}

// inserts an ad row
func CreateAdvertisement(userId, levelID string, adType int, imageURL string) (int64, error) {
	if userId == "" || levelID == "" || imageURL == "" {
		return 0, fmt.Errorf("missing ad fields")
	}

	stmt, err := prepareStmt(data, "INSERT INTO advertisements (user_id, level_id, type, image_url) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(userId, levelID, adType, imageURL)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

// fetches all ads for a given user
func ListAllAdvertisements() ([]AdRow, error) {
	stmt, err := prepareStmt(data, "SELECT ad_id, user_id, level_id, type, image_url, created_at FROM advertisements ORDER BY ad_id DESC")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
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

func FilterAdsByUser(rows []AdRow, userId string) ([]AdRow, error) {
	var out []AdRow
	for _, r := range rows {
		if r.UserID == userId {
			out = append(out, r)
		}
	}

	return out, nil
}

func FilterAdsByType(rows []AdRow, adType AdType) ([]AdRow, error) {
	// Map type to number
	typeNum := 0
	switch adType {
	case AdTypeBanner:
		typeNum = 1

	case AdTypeSquare:
		typeNum = 2

	case AdTypeSkyscraper:
		typeNum = 3

	default:
		return nil, fmt.Errorf("invalid ad type")
	}

	var out []AdRow
	for _, r := range rows {
		if r.Type == typeNum {
			out = append(out, r)
		}
	}

	return out, nil
}

// returns the owning user_id for an ad
func GetAdvertisementOwner(adId int64) (string, error) {
	var uid string

	stmt, err := prepareStmt(data, "SELECT user_id FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return "", err
	}

	err = stmt.QueryRow(adId).Scan(&uid)
	if err != nil {
		return "", err
	}

	return uid, nil
}

func DeleteAdvertisement(adId int64) error {
	stmt, err := prepareStmt(data, "DELETE FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(adId)
	if err != nil {
		return err
	}

	return nil
}

func DeleteAllExpiredAds() error {
	stmt, err := prepareStmt(data, "DELETE FROM advertisements WHERE created_at < NOW() - INTERVAL 7 DAY")
	if err != nil {
		return err
	}

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	return nil
}

func init() {
	var err error

	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	log.Info("Connecting to database with URI: " + uri)
	data, err = sql.Open("mysql", uri)
	if err != nil {
		log.Error("Failed to establish MariaDB connection: %s", err.Error())
		return
	}

	err = data.Ping()
	if err != nil {
		log.Error("Failed to ping database: %s", err.Error())
		return
	} else if data == nil {
		log.Error("Database connection is nil")
		return
	}

	log.Print("MariaDB connection established.")
}

// returns total_views and total_clicks for a given user id
func GetUserTotals(userId string) (int, int, error) {
	if userId == "" {
		return 0, 0, fmt.Errorf("empty user id")
	}

	stmt, err := prepareStmt(data, "SELECT total_views, total_clicks FROM users WHERE id = ?")
	if err != nil {
		return 0, 0, err
	}

	var views int
	var clicks int
	err = stmt.QueryRow(userId).Scan(&views, &clicks)
	if err != nil {
		return 0, 0, err
	}

	return views, clicks, nil
}
