package database

import (
	"database/sql"
	"fmt"

	"service/log"
	"service/utils"
)

var dat *sql.DB

// Register a new client event for an ad
func NewStat(event utils.AdEvent, adId int64, userID string) error {
	log.Debug("Registering new %s for user %s on ad %d", event, userID, adId)

	query := fmt.Sprintf("UPDATE advertisements SET %s = %s + 1 WHERE ad_id = ?", event, event)

	stmt, err := utils.PrepareStmt(dat, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	if _, err := stmt.Exec(adId); err != nil {
		return err
	}

	_, err = GetAdvertisement(adId)
	if err != nil {
		log.Error("Failed to get advertisement %d: %s", adId, err.Error())
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

	// Get count of active (non-pending) advertisements
	err := dat.QueryRow("SELECT COUNT(*) FROM advertisements WHERE pending = FALSE").Scan(&adCount)
	if err != nil {
		log.Error("Failed to fetch ad count: %s", err.Error())
		return 0, 0, 0, err
	}

	return totalViews, totalClicks, adCount, nil
}

func init() {
	dat = utils.Db()
}
