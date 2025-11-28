package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"service/log"
	"service/utils"

	"github.com/patrickmn/go-cache"
)

var dat *sql.DB
var globals = cache.New(5*time.Minute, 10*time.Minute)

// Register a new client event for an ad
func NewStat(event utils.AdEvent, adId int64) error {
	log.Debug("Registering new %s on ad %d", event, adId)

	query := fmt.Sprintf("UPDATE advertisements SET %s = %s + 1 WHERE ad_id = ?", event, event)

	stmt, err := utils.PrepareStmt(dat, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(adId); err != nil {
		return err
	}

	ad, err := GetAdvertisement(adId)
	if err != nil {
		log.Error("Failed to get advertisement %d: %s", adId, err.Error())
		return err
	}

	var viewsDelta, clicksDelta int = 0, 0
	switch event {
	case utils.AdEventView:
		viewsDelta = 1
		ad.Views += uint64(viewsDelta)
	case utils.AdEventClick:
		clicksDelta = 1
		ad.Clicks += uint64(clicksDelta)
	default:
		return fmt.Errorf("invalid ad event")
	}

	// Get the ad owner and increment their stats
	if ownerID, ownerErr := GetAdvertisementOwnerId(adId); ownerErr == nil && ownerID != "" {
		log.Debug("Incrementing stats for owner %s: views +%d, clicks +%d", ownerID, viewsDelta, clicksDelta)
		if incErr := IncrementUserStats(ownerID, viewsDelta, clicksDelta); incErr != nil {
			log.Error("Failed to increment total stats for user %s: %s", ownerID, incErr.Error())
		}
	} else {
		log.Warn("Could not find owner for ad %d: %v", adId, ownerErr)
	}

	log.Debug("Successfully registered %s for ad %d", event, adId)
	return nil
}

// returns total_views and total_clicks for a given user id
func GetUserTotals(userId string) (utils.Stats, error) {
	if val, found := globals.Get(userId); found {
		log.Debug("Returning cached global stats for user of ID %s", userId)
		return val.(utils.Stats), nil
	}

	var stats = utils.Stats{}

	if userId == "" {
		return stats, fmt.Errorf("empty user id")
	}

	stmt, err := utils.PrepareStmt(dat, "SELECT total_views, total_clicks FROM users WHERE id = ?")
	if err != nil {
		return stats, err
	}
	defer stmt.Close()

	err = stmt.QueryRow(userId).Scan(&stats.Views, &stats.Clicks)
	if err != nil {
		return stats, err
	}

	globals.Set(userId, stats, cache.DefaultExpiration)

	return stats, nil
}

// GetGlobalStats returns the total views, total clicks, and count of active advertisements
func GetGlobalStats() (utils.GlobalStats, error) {
	if val, found := globals.Get("global"); found {
		log.Debug("Returning cached global stats")
		return val.(utils.GlobalStats), nil
	}

	stats := utils.GlobalStats{}

	countStmt, err := utils.PrepareStmt(dat, "SELECT total_views, total_clicks FROM users WHERE banned = FALSE")
	if err != nil {
		return stats, err
	}
	defer countStmt.Close()

	countRows, err := countStmt.Query()
	if err != nil {
		return stats, err
	}
	defer countRows.Close()

	for countRows.Next() {
		var cr utils.Ad
		if err := countRows.Scan(
			&cr.Views,
			&cr.Clicks,
		); err != nil {
			log.Error("Failed to scan row for global stats: %s", err.Error())
		}

		stats.TotalViews += cr.Views
		stats.TotalClicks += cr.Clicks
	}

	adStmt, err := utils.PrepareStmt(dat, "SELECT COUNT(*) FROM advertisements WHERE pending = FALSE")
	if err != nil {
		return stats, err
	}
	defer adStmt.Close()

	err = adStmt.QueryRow().Scan(&stats.AdCount)
	if err != nil {
		return stats, err
	}

	globals.Set("global", stats, cache.DefaultExpiration)

	return stats, nil
}

func GetModDownloads() (uint64, error) {
	req, err := http.NewRequest(http.MethodGet, "https://api.geode-sdk.org/v1/mods/arcticwoof.player_advertisements", nil)
	if err != nil {
		return 0, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	dlBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	log.Debug("Mod endpoint status: %s", resp.Status)

	if !strings.Contains(resp.Header.Get("Content-Type"), "application/json") {
		log.Error("Geode returned non-JSON: %s", string(dlBody))
		return 0, fmt.Errorf("request returned non-json")
	}

	if resp.Request != nil {
		log.Debug("Mod endpoint final URL: %s", resp.Request.URL.String())
	}

	type payload struct {
		DownloadCount uint64 `json:"download_count"`
	}

	var dlResp struct {
		Error   string  `json:"error"`
		Payload payload `json:"payload"`
	}

	if err := json.Unmarshal(dlBody, &dlResp); err != nil {
		log.Error("Failed to decode mod response: %s", err.Error())
		return 0, err
	}

	if dlResp.Error != "" {
		return 0, fmt.Errorf("error: %s", dlResp.Error)
	}

	return dlResp.Payload.DownloadCount, nil
}

func init() {
	dat = utils.Db()
}
