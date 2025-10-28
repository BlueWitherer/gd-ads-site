package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"
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
	StatByViews  StatBy = "total_views"  // Filter stats by ad
	StatByClicks StatBy = "total_clicks" // Filter stats by user
)

type User struct {
	ID          string `json:"id"`           // Discord user ID
	Username    string `json:"username"`     // Discord username
	TotalViews  int    `json:"total_views"`  // Total registered views on all ads
	TotalClicks int    `json:"total_clicks"` // Total registered clicks on all ads
	IsAdmin     bool   `json:"is_admin"`     // Active administrator status
	IsStaff     bool   `json:"is_staff"`     // Active staff status
	Verified    bool   `json:"verified"`     // Trusted status
	Banned      bool   `json:"banned"`       // Banned status
	Created     string `json:"created_at"`   // First created
	Updated     string `json:"updated_at"`   // Last updated
}

// Database row for advertisements listing
type Ad struct {
	AdID       int64  `json:"ad_id"`                 // Advertisement ID
	UserID     string `json:"user_id"`               // Owner Discord user ID
	LevelID    int64  `json:"level_id"`              // Geometry Dash level ID
	Type       int    `json:"type"`                  // Type of advertisement
	ViewCount  int    `json:"view_count,omitempty"`  // Total registered views
	ClickCount int    `json:"click_count,omitempty"` // Total registered clicks
	ImageURL   string `json:"image_url"`             // URL to the advertisement image
	Created    string `json:"created_at"`            // First created
	Expiry     int64  `json:"expiry"`                // Date of expiration
	Pending    bool   `json:"pending"`               // Under review
}

// private method to safely prepare the sql statement
func prepareStmt(db *sql.DB, sql string) (*sql.Stmt, error) {
	if db != nil {
		log.Debug("Preparing connection for statement %s", sql)
		return db.Prepare(sql)
	} else {
		return nil, fmt.Errorf("database connection non-existent")
	}
}

func AdTypeFromInt(t int) (AdType, error) {
	switch t {
	case 1:
		return AdTypeBanner, nil
	case 2:
		return AdTypeSquare, nil
	case 3:
		return AdTypeSkyscraper, nil

	default:
		return "", fmt.Errorf("invalid ad type")
	}
}

func IntFromAdType(t AdType) (int, error) {
	switch t {
	case AdTypeBanner:
		return 1, nil
	case AdTypeSquare:
		return 2, nil
	case AdTypeSkyscraper:
		return 3, nil

	default:
		return 0, fmt.Errorf("invalid ad type")
	}
}

// Register a new client event for an ad
func NewStat(event AdEvent, adId int64, user int64) error {
	log.Debug("Registering new %s", event)
	sql := fmt.Sprintf("INSERT INTO %s (ad_id, user_id, timestamp) VALUES (?, ?, ?)", event)

	stmt, err := prepareStmt(data, sql)
	if err != nil {
		return err
	}

	var viewsDelta, clicksDelta int = 0, 0
	switch event {
	case AdEventView:
		viewsDelta = 1
	case AdEventClick:
		clicksDelta = 1

	default:
		return fmt.Errorf("invalid ad event")
	}

	if ownerID, ownerErr := GetAdvertisementOwnerId(adId); ownerErr == nil && ownerID != "" {
		if incErr := IncrementUserStats(ownerID, viewsDelta, clicksDelta); incErr != nil {
			log.Error("Failed to increment total clicks: %s", incErr.Error())
		}
	}

	_, err = stmt.Exec(adId, user, time.Now())
	return err
}

func NewStatWithUserID(event AdEvent, adId int64, userID string) error {
	log.Debug("Registering new %s for user %s", event, userID)
	sql := fmt.Sprintf("INSERT INTO %s (ad_id, user_id, timestamp) VALUES (?, ?, ?)", event)

	stmt, err := prepareStmt(data, sql)
	if err != nil {
		return err
	}

	var viewsDelta, clicksDelta int = 0, 0
	switch event {
	case AdEventView:
		viewsDelta = 1
	case AdEventClick:
		clicksDelta = 1

	default:
		return fmt.Errorf("invalid ad event")
	}

	if ownerID, ownerErr := GetAdvertisementOwnerId(adId); ownerErr == nil && ownerID != "" {
		if incErr := IncrementUserStats(ownerID, viewsDelta, clicksDelta); incErr != nil {
			log.Error("Failed to increment total stats: %s", incErr.Error())
		}
	}

	_, err = stmt.Exec(adId, userID, time.Now())
	return err
}

func GetUser(id string) (User, error) {
	if id == "" {
		return User{}, fmt.Errorf("empty user id")
	}

	stmt, err := prepareStmt(data, "SELECT username, id, total_views, total_clicks, is_admin, banned, created_at, updated_at FROM users WHERE id = ?")
	if err != nil {
		return User{}, err
	}

	var user User
	err = stmt.QueryRow(id).Scan(&user.Username, &user.ID, &user.TotalViews, &user.TotalClicks, &user.IsAdmin, &user.Banned, &user.Created, &user.Updated)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func GetAllUsers() ([]User, error) {
	stmt, err := prepareStmt(data, "SELECT id, username, total_clicks, total_views, is_admin, banned, created_at, updated_at FROM users ORDER BY id DESC")
	if err != nil {
		return nil, err
	}

	users, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer users.Close()

	var out []User
	for users.Next() {
		var u User
		if err := users.Scan(&u.ID, &u.Username, &u.TotalClicks, &u.TotalViews, &u.IsAdmin, &u.Banned, &u.Created, &u.Updated); err != nil {
			return nil, err
		}

		out = append(out, u)
	}

	return out, users.Err()
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

// inserts or updates an ad row
func CreateAdvertisement(userId string, levelID string, adType int, imageURL string) (int64, error) {
	if userId == "" || levelID == "" || imageURL == "" {
		return 0, fmt.Errorf("missing ad fields")
	}

	// Create new ad - allow multiple ads per user per type
	stmt, err := prepareStmt(data, "INSERT INTO advertisements (user_id, level_id, type, image_url, pending) VALUES (?, ?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(userId, levelID, adType, imageURL, true)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func ApproveAd(id int64) (Ad, error) {
	stmt, err := prepareStmt(data, "UPDATE advertisements SET pending = FALSE WHERE ad_id = ?")
	if err != nil {
		return Ad{}, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return Ad{}, err
	}

	return GetAdvertisement(id)
}

func BanUser(id string) (User, error) {
	// delete all advertisements associated with the user
	deleteAdsStmt, err := prepareStmt(data, "DELETE FROM advertisements WHERE user_id = ?")
	if err != nil {
		return User{}, err
	}

	_, err = deleteAdsStmt.Exec(id)
	if err != nil {
		return User{}, err
	}

	// ban the user
	stmt, err := prepareStmt(data, "UPDATE users SET banned = TRUE WHERE id = ?")
	if err != nil {
		return User{}, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return User{}, err
	}

	return GetUser(id)
}

func UnbanUser(id string) (User, error) {
	// unban the user
	stmt, err := prepareStmt(data, "UPDATE users SET banned = FALSE WHERE id = ?")
	if err != nil {
		return User{}, err
	}

	_, err = stmt.Exec(id)
	if err != nil {
		return User{}, err
	}

	return GetUser(id)
}

func GetAdUnixExpiry(ad Ad) (int64, error) {
	t, err := time.Parse("2006-01-02 15:04:05", ad.Created)
	if err != nil {
		return 0, err
	}

	expiry := t.Unix() + int64((7 * 24 * time.Hour).Seconds())

	return expiry, err
}

// fetches all ads for a given user
func ListAllAdvertisements() ([]Ad, error) {
	stmt, err := prepareStmt(data, "SELECT ad_id, user_id, level_id, type, image_url, created_at, pending FROM advertisements ORDER BY ad_id DESC")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var out []Ad
	for rows.Next() {
		var r Ad
		if err := rows.Scan(&r.AdID, &r.UserID, &r.LevelID, &r.Type, &r.ImageURL, &r.Created, &r.Pending); err != nil {
			return nil, err
		}

		r.Expiry, err = GetAdUnixExpiry(r)
		if err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	return out, rows.Err()
}

func ListPendingAdvertisements() ([]Ad, error) {
	// Use != 0 to match tinyint(1) values in MySQL/MariaDB
	stmt, err := prepareStmt(data, "SELECT ad_id, user_id, level_id, type, image_url, created_at, pending FROM advertisements WHERE pending = TRUE ORDER BY ad_id DESC")
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	out := make([]Ad, 0)
	for rows.Next() {
		var r Ad
		if err := rows.Scan(&r.AdID, &r.UserID, &r.LevelID, &r.Type, &r.ImageURL, &r.Created, &r.Pending); err != nil {
			return nil, err
		}

		r.Expiry, err = GetAdUnixExpiry(r)
		if err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	return out, rows.Err()
}

func FilterAdsByPending(rows []Ad, showPending bool) ([]Ad, error) {
	out := make([]Ad, 0)
	for _, r := range rows {
		if r.Pending == showPending {
			out = append(out, r)
		}
	}

	return out, nil
}

func FilterAdsFromBannedUsers(rows []Ad) ([]Ad, error) {
	var out []Ad
	for _, r := range rows {
		user, err := GetUser(r.UserID)
		if err != nil {
			return nil, err
		}

		if !user.Banned {
			out = append(out, r)
		}
	}

	return out, nil
}

func FilterAdsByUser(rows []Ad, userId string) ([]Ad, error) {
	var out []Ad
	for _, r := range rows {
		if r.UserID == userId {
			out = append(out, r)
		}
	}

	return out, nil
}

func FilterAdsByType(rows []Ad, adType AdType) ([]Ad, error) {
	typeNum, err := IntFromAdType(adType)
	if err != nil {
		return nil, err
	}

	var out []Ad
	for _, r := range rows {
		if r.Type == typeNum {
			out = append(out, r)
		}
	}

	return out, nil
}

func UserLeaderboard(stat StatBy, page uint64, maxPerPage uint64) ([]User, error) {
	stmt, err := prepareStmt(data, fmt.Sprintf("SELECT * FROM users WHERE banned = FALSE ORDER BY %s DESC", stat))
	if err != nil {
		return nil, err
	}

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	var out []User
	for rows.Next() {
		var r User

		if err := rows.Scan(&r.ID, &r.Username, &r.TotalViews, &r.TotalClicks, &r.IsAdmin, &r.Banned, &r.Created, &r.Updated); err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	start := page * maxPerPage
	end := start + maxPerPage

	if start >= uint64(len(out)) {
		return []User{}, nil
	}

	if end > uint64(len(out)) {
		end = uint64(len(out))
	}

	return out[start:end], nil
}

// get an advertisement by id
func GetAdvertisement(adId int64) (Ad, error) {
	stmt, err := prepareStmt(data, "SELECT ad_id, user_id, level_id, type, image_url, created_at, pending FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return Ad{}, err
	}

	// QueryRow is more convenient when expecting a single row
	row := stmt.QueryRow(adId)
	if row != nil {
		var r Ad
		if err := row.Scan(&r.AdID, &r.UserID, &r.LevelID, &r.Type, &r.ImageURL, &r.Created, &r.Pending); err != nil {
			if err == sql.ErrNoRows {
				return Ad{}, nil
			}

			return Ad{}, err
		}

		r.Expiry, err = GetAdUnixExpiry(r)
		if err != nil {
			return Ad{}, err
		}

		return r, nil
	} else {
		return Ad{}, fmt.Errorf("ad not found")
	}
}

// returns the owning user_id for an ad
func GetAdvertisementOwnerId(adId int64) (string, error) {
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

func UpdateAdvertisementImageURL(adId int64, imageURL string) error {
	if imageURL == "" {
		return fmt.Errorf("empty image url")
	}

	stmt, err := prepareStmt(data, "UPDATE advertisements SET image_url = ? WHERE ad_id = ?")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(imageURL, adId)
	return err
}

func DeleteAdvertisement(adId int64) (Ad, error) {
	ad, err := GetAdvertisement(adId)
	if err != nil {
		return ad, err
	}

	stmt, err := prepareStmt(data, "DELETE FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return ad, err
	}

	_, err = stmt.Exec(adId)
	if err != nil {
		return ad, err
	}

	return ad, nil
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

// CountActiveAdvertisementsByUser returns the count of active (non-expired) advertisements for a user
func CountActiveAdvertisementsByUser(userId string) (int, error) {
	if userId == "" {
		return 0, fmt.Errorf("empty user id")
	}

	stmt, err := prepareStmt(data, "SELECT COUNT(*) FROM advertisements WHERE user_id = ? AND created_at > NOW() - INTERVAL 7 DAY")
	if err != nil {
		return 0, err
	}

	var count int
	err = stmt.QueryRow(userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
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

// returns total_views and total_clicks for a given ad id
func GetAdStats(adId int64) (int, int, error) {
	if adId == 0 {
		return 0, 0, fmt.Errorf("invalid ad id")
	}

	// Count views for this ad
	viewStmt, err := prepareStmt(data, "SELECT COUNT(*) FROM views WHERE ad_id = ?")
	if err != nil {
		return 0, 0, err
	}

	var views int
	err = viewStmt.QueryRow(adId).Scan(&views)
	if err != nil {
		return 0, 0, err
	}

	// Count clicks for this ad
	clickStmt, err := prepareStmt(data, "SELECT COUNT(*) FROM clicks WHERE ad_id = ?")
	if err != nil {
		return 0, 0, err
	}

	var clicks int
	err = clickStmt.QueryRow(adId).Scan(&clicks)
	if err != nil {
		return 0, 0, err
	}

	return views, clicks, nil
}

// initializeSchema reads and executes the schema.sql file to create tables if they don't exist
func initializeSchema() error {
	schemaPath := filepath.Join("service", "database", "schema.sql")
	log.Debug("Reading database schema from %s", schemaPath)

	schemaSQL, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("failed to read schema file: %w", err)
	}

	// Split the SQL file into individual statements
	statements := strings.Split(string(schemaSQL), ";")

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" || strings.HasPrefix(stmt, "--") {
			continue
		}

		log.Debug("Executing schema statement: %.50s...", stmt)
		_, err := data.Exec(stmt)
		if err != nil {
			return fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	log.Done("Database schema initialized successfully")
	return nil
}

// GetGlobalStats returns the total views, total clicks, and count of active advertisements
func GetGlobalStats() (int, int, int, error) {
	if data == nil {
		return 0, 0, 0, fmt.Errorf("database connection non-existent")
	}

	var totalViews, totalClicks, adCount int

	// Get sum of total_views and total_clicks from users table
	err := data.QueryRow("SELECT COALESCE(SUM(total_views), 0), COALESCE(SUM(total_clicks), 0) FROM users").Scan(&totalViews, &totalClicks)
	if err != nil {
		log.Error("Failed to fetch user stats: %s", err.Error())
		return 0, 0, 0, err
	}

	// Get count of active (non-pending) advertisements
	err = data.QueryRow("SELECT COUNT(*) FROM advertisements WHERE pending = FALSE").Scan(&adCount)
	if err != nil {
		log.Error("Failed to fetch ad count: %s", err.Error())
		return 0, 0, 0, err
	}

	return totalViews, totalClicks, adCount, nil
}

func init() {
	var err error

	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASS"),
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
	)

	log.Info("Connecting to database with URI: %s", uri)
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

	// Initialize database schema (create tables if they don't exist)
	if err := initializeSchema(); err != nil {
		log.Error("Failed to initialize database schema: %s", err.Error())
		return
	}
}
