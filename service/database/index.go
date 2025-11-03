package database

import (
	"database/sql"
	"fmt"
	"time"

	"service/log"
	"service/utils"
)

var dat *sql.DB

// Register a new client event for an ad
func NewStat(event utils.AdEvent, adId int64, user int64) error {
	log.Debug("Registering new %s on ad %d for user %d", event, adId, user)
	sql := fmt.Sprintf("INSERT INTO %s (ad_id, user_id, timestamp) VALUES (?, ?, ?)", event)

	stmt, err := utils.PrepareStmt(dat, sql)
	if err != nil {
		return err
	}

	var viewsDelta, clicksDelta int = 0, 0
	switch event {
	case utils.AdEventView:
		viewsDelta = 1
	case utils.AdEventClick:
		clicksDelta = 1

	default:
		return fmt.Errorf("invalid ad event")
	}

	if ownerID, ownerErr := GetAdvertisementOwnerId(adId); ownerErr == nil && ownerID != "" {
		log.Info("Incrementing stats for owner %s: views +%d, clicks +%d", ownerID, viewsDelta, clicksDelta)
		if incErr := IncrementUserStats(ownerID, viewsDelta, clicksDelta); incErr != nil {
			log.Error("Failed to increment total clicks: %s", incErr.Error())
		}
	} else {
		log.Warn("Could not find owner for ad %d: %v", adId, ownerErr)
	}

	_, err = stmt.Exec(adId, user, time.Now())
	if err != nil {
		log.Error("Failed to insert %s record: %s", event, err.Error())
		return err
	}

	log.Debug("Successfully registered %s for ad %d", event, adId)
	return err
}

func NewStatWithUserID(event utils.AdEvent, adId int64, userID string) error {
	log.Debug("Registering new %s for user %s on ad %d", event, userID, adId)

	_, err := GetAdvertisement(adId)
	if err != nil {
		log.Error("Failed to get advertisement %d: %s", adId, err.Error())
		return err
	}

	sql := fmt.Sprintf("INSERT INTO %s (ad_id, user_id, timestamp) VALUES (?, ?, ?)", event)

	stmt, err := utils.PrepareStmt(dat, sql)
	if err != nil {
		return err
	}

	var viewsDelta, clicksDelta int = 0, 0
	switch event {
	case utils.AdEventView:
		viewsDelta = 1
	case utils.AdEventClick:
		clicksDelta = 1

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

	_, err = stmt.Exec(adId, userID, time.Now())
	if err != nil {
		log.Error("Failed to insert %s record: %s", event, err.Error())
		return err
	}

	log.Debug("Successfully registered %s for ad %d", event, adId)
	return nil
}

// returns total_views and total_clicks for a given user id
func GetUserTotals(userId string) (int, int, error) {
	if userId == "" {
		return 0, 0, fmt.Errorf("empty user id")
	}

	stmt, err := utils.PrepareStmt(dat, "SELECT total_views, total_clicks FROM users WHERE id = ?")
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

// GetGlobalStats returns the total views, total clicks, and count of active advertisements
func GetGlobalStats() (uint64, uint64, uint, error) {
	if dat == nil {
		return 0, 0, 0, fmt.Errorf("database connection non-existent")
	}

	var totalViews, totalClicks uint64
	var adCount uint

	err := dat.QueryRow("SELECT COUNT(*) FROM views").Scan(&totalViews)
	if err != nil {
		log.Error("Failed to fetch view stats: %s", err.Error())
		return 0, 0, 0, err
	}

	err = dat.QueryRow("SELECT COUNT(*) FROM clicks").Scan(&totalClicks)
	if err != nil {
		log.Error("Failed to fetch click stats: %s", err.Error())
		return 0, 0, 0, err
	}

	// Get count of active (non-pending) advertisements
	err = dat.QueryRow("SELECT COUNT(*) FROM advertisements WHERE pending = FALSE").Scan(&adCount)
	if err != nil {
		log.Error("Failed to fetch ad count: %s", err.Error())
		return 0, 0, 0, err
	}

	return totalViews, totalClicks, adCount, nil
}

func init() {
	dat = utils.Db()
}
