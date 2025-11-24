package database

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"service/log"
	"service/utils"
)

func newAds() *[]*utils.Ad {
	return new([]*utils.Ad)
}

// Current ads cache
var currentAds *[]*utils.Ad = nil
var currentAdsSince time.Time = time.Now()

func getAds() *[]*utils.Ad {
	if currentAds != nil {
		log.Debug("Returning cached ads list")
		return currentAds
	}

	currentAdsSince = time.Now()

	return newAds()
}

func findAd(id int64) (*utils.Ad, bool) {
	if currentAds != nil {
		for _, a := range *currentAds {
			if a.AdID == id {
				return a, true
			}
		}
	}

	return nil, false
}

func setAd(ad *utils.Ad) *[]*utils.Ad {
	if currentAds != nil {
		log.Debug("Caching ad %d", ad.AdID)
		for i, a := range *currentAds {
			if a.AdID == ad.AdID {
				(*currentAds)[i] = ad
				return getAds()
			}
		}

		*currentAds = append(*currentAds, ad)
	}

	return getAds()
}

func deleteAd(id int64) *[]*utils.Ad {
	if currentAds != nil {
		for i, a := range *currentAds {
			if a.AdID == id {
				*currentAds = append((*currentAds)[:i], (*currentAds)[i+1:]...)
			}
		}
	}

	return getAds()
}

func ApproveAd(id int64) (*utils.Ad, error) {
	stmt, err := utils.PrepareStmt(dat, "UPDATE advertisements SET pending = FALSE, created_at = NOW() WHERE ad_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return nil, err
	}

	// fetch the ad so we can return it and touch its image file
	ad, err := GetAdvertisement(id)
	if err != nil {
		return nil, err
	}

	// try to touch the ad image to reset its modification time
	if ad != nil {
		if adType, err := utils.AdTypeFromInt(ad.Type); err == nil {
			adPath := filepath.Join("..", "ad_storage", string(adType), fmt.Sprintf("%s-%d.webp", ad.UserID, ad.AdID))
			now := time.Now()

			if err := os.Chtimes(adPath, now, now); err != nil {
				log.Error("Failed to reset image for ad approval %s: %s", adPath, err.Error())
			} else {
				log.Info("Reset image %s for ad approval", adPath)
			}
		} else {
			log.Error("Failed to determine ad type for resetting file: %s", err.Error())
		}

		currentAds = setAd(ad)
	}

	return ad, nil
}

// inserts or updates an ad row
func CreateAdvertisement(userId string, levelID string, adType int) (int64, error) {
	if userId == "" || levelID == "" {
		return 0, fmt.Errorf("missing ad fields")
	}

	// Create new ad - allow multiple ads per user per type
	stmt, err := utils.PrepareStmt(dat, "INSERT INTO advertisements (user_id, level_id, type, pending) VALUES (?, ?, ?, ?)")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	res, err := stmt.Exec(userId, levelID, adType, true)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func GetAdUnixExpiry(ad *utils.Ad) int64 {
	expiry := ad.Created.Unix() + int64((7 * 24 * time.Hour).Seconds())

	return expiry
}

// fetches all ads for a given user
func ListAllAdvertisements() ([]*utils.Ad, error) {
	if time.Since(currentAdsSince) > 15*time.Minute {
		currentAds = nil
	}

	if currentAds != nil && len(*currentAds) > 0 {
		log.Debug("Returning cached ads list")
		return *getAds(), nil
	}

	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM advertisements ORDER BY ad_id DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*utils.Ad
	for rows.Next() {
		r := new(utils.Ad)
		if err := rows.Scan(
			&r.AdID,
			&r.UserID,
			&r.LevelID,
			&r.Type,
			&r.Views,
			&r.Clicks,
			&r.ImageURL,
			&r.Created,
			&r.Pending,
			&r.BoostCount,
		); err != nil {
			return nil, err
		}

		r.Expiry = GetAdUnixExpiry(r)

		currentAds = setAd(r)

		out = append(out, r)
	}

	return out, rows.Err()
}

func ListPendingAdvertisements() ([]*utils.Ad, error) {
	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM advertisements WHERE pending = TRUE ORDER BY ad_id DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	out := make([]*utils.Ad, 0)
	for rows.Next() {
		r := new(utils.Ad)
		if err := rows.Scan(
			&r.AdID,
			&r.UserID,
			&r.LevelID,
			&r.Type,
			&r.Views,
			&r.Clicks,
			&r.ImageURL,
			&r.Created,
			&r.Pending,
			&r.BoostCount,
		); err != nil {
			return nil, err
		}

		r.Expiry = GetAdUnixExpiry(r)

		currentAds = setAd(r)

		out = append(out, r)
	}

	return out, rows.Err()
}

func FilterAdsByPending(rows []*utils.Ad, showPending bool) ([]*utils.Ad, error) {
	out := make([]*utils.Ad, 0)
	for _, r := range rows {
		if r.Pending == showPending {
			out = append(out, r)
		}
	}

	return out, nil
}

func FilterAdsFromBannedUsers(rows []*utils.Ad) ([]*utils.Ad, error) {
	var out []*utils.Ad
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

func FilterAdsByUser(rows []*utils.Ad, userId string) ([]*utils.Ad, error) {
	var out []*utils.Ad
	for _, r := range rows {
		if r.UserID == userId {
			out = append(out, r)
		}
	}

	return out, nil
}

func FilterAdsByType(rows []*utils.Ad, adType utils.AdType) ([]*utils.Ad, error) {
	typeNum, err := utils.IntFromAdType(adType)
	if err != nil {
		return nil, err
	}

	var out []*utils.Ad
	for _, r := range rows {
		if r.Type == typeNum {
			out = append(out, r)
		}
	}

	return out, nil
}

func GetAdvertisement(adId int64) (*utils.Ad, error) {
	if val, found := findAd(adId); found {
		views, clicks, err := GetAdStats(adId)
		if err != nil {
			return nil, err
		}

		val.Views = uint64(views)
		val.Clicks = uint64(clicks)

		return val, nil
	}

	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	row := stmt.QueryRow(adId)
	if row != nil {
		r := new(utils.Ad)
		if err := row.Scan(
			&r.AdID,
			&r.UserID,
			&r.LevelID,
			&r.Type,
			&r.Views,
			&r.Clicks,
			&r.ImageURL,
			&r.Created,
			&r.Pending,
			&r.BoostCount,
		); err != nil {
			if err == sql.ErrNoRows {
				return nil, err
			}

			return nil, err
		}

		r.Expiry = GetAdUnixExpiry(r)

		currentAds = setAd(r)

		return r, nil
	} else {
		return nil, fmt.Errorf("ad not found")
	}
}

// returns the owning user_id for an ad
func GetAdvertisementOwnerId(adId int64) (string, error) {
	if val, found := findAd(adId); found {
		return val.UserID, nil
	}

	var uid string

	stmt, err := utils.PrepareStmt(dat, "SELECT user_id FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

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

	stmt, err := utils.PrepareStmt(dat, "UPDATE advertisements SET image_url = ? WHERE ad_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	ad, err := GetAdvertisement(adId)
	if err != nil {
		return err
	}

	ad.ImageURL = imageURL
	currentAds = setAd(ad)

	_, err = stmt.Exec(imageURL, adId)
	return err
}

func DeleteAdvertisement(adId int64) (*utils.Ad, error) {
	ad, err := GetAdvertisement(adId)
	if err != nil {
		return ad, err
	}

	stmt, err := utils.PrepareStmt(dat, "DELETE FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return ad, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(adId)
	if err != nil {
		return ad, err
	}

	adType, err := utils.AdTypeFromInt(ad.Type)
	if err != nil {
		return ad, err
	}

	adDir := filepath.Join("..", "ad_storage", string(adType), fmt.Sprintf("%s-%d.webp", ad.UserID, ad.AdID))
	err = os.Remove(adDir)
	if err != nil {
		return ad, err
	}

	currentAds = deleteAd(adId)

	return ad, nil
}

func DeleteAllExpiredAds() error {
	stmt, err := utils.PrepareStmt(dat, "DELETE FROM advertisements WHERE created_at < NOW() - INTERVAL 7 DAY")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		return err
	}

	adsDir := filepath.Join("..", "ad_storage")
	err = filepath.WalkDir(adsDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			log.Error("Error accessing path %s: %s", path, err.Error())
			return nil // continue walking
		}

		if d.IsDir() {
			return nil // skip directories
		}

		info, err := d.Info()
		if err != nil {
			log.Error("Failed to get file info for %s: %s", path, err.Error())
			return nil
		}

		if time.Since(info.ModTime()) > 7*24*time.Hour {
			log.Info("Removing expired ad %s (%v B)", path, info.Size())
			if err := os.Remove(path); err != nil {
				log.Error("Failed to remove file %s: %s", path, err.Error())
			}
		} else {
			log.Debug("Advertisement %s is still valid", path)
		}

		return nil
	})

	if err != nil {
		log.Error("Failed to walk ad directory: %s", err.Error())
		return err
	}

	currentAds = nil // clear cache

	return nil
}

// returns the count of active (non-expired) advertisements for a user
func CountActiveAdvertisementsByUser(userId string) (int, error) {
	if userId == "" {
		return 0, fmt.Errorf("empty user id")
	}

	stmt, err := utils.PrepareStmt(dat, "SELECT COUNT(*) FROM advertisements WHERE user_id = ? AND created_at > NOW() - INTERVAL 7 DAY")
	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var count int
	err = stmt.QueryRow(userId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

// returns total_views and total_clicks for a given ad id
func GetAdStats(adId int64) (int, int, error) {
	stmt, err := utils.PrepareStmt(dat, "SELECT views, clicks FROM advertisements WHERE ad_id = ?")
	if err != nil {
		return 0, 0, err
	}
	defer stmt.Close()

	var views int
	var clicks int
	err = stmt.QueryRow(adId).Scan(&views, &clicks)
	if err != nil {
		return 0, 0, err
	}

	return views, clicks, nil
}

func BoostAd(adId int64, boosts uint, user string) error {
	deductStmt, err := utils.PrepareStmt(dat, "UPDATE users SET boost_count = boost_count - ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer deductStmt.Close()

	_, err = deductStmt.Exec(boosts, user)
	if err != nil {
		return err
	}

	stmt, err := utils.PrepareStmt(dat, "UPDATE advertisements SET boost_count = boost_count + ? WHERE ad_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(adId, boosts)
	if err != nil {
		return err
	}

	return nil
}

func AddBoostsToUser(userId string, boosts uint) error {
	stmt, err := utils.PrepareStmt(dat, "UPDATE users SET boost_count = boost_count + ? WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(boosts, userId)
	if err != nil {
		return err
	}

	return nil
}

func NewReport(adId int64, accountId int, description string) error {
	stmt, err := utils.PrepareStmt(dat, "INSERT INTO reports (ad_id, account_id, description) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(adId, accountId, description)
	if err != nil {
		return err
	}

	return nil
}

func GetReport(id int64) (*utils.Report, error) {
	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM reports WHERE id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	report := new(utils.Report)
	var adId int64
	err = stmt.QueryRow(id).Scan(&report.ID, &adId, &report.AccountID, &report.Description, &report.Created)
	if err != nil {
		return nil, err
	}

	ad, err := GetAdvertisement(adId)
	if err != nil {
		return nil, err
	}

	report.Ad = *ad

	return report, nil
}

func ListAllReports() ([]*utils.Report, error) {
	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM reports ORDER BY created_at ASC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []*utils.Report
	for rows.Next() {
		r := new(utils.Report)
		var adId int64
		if err := rows.Scan(
			&adId,
			&r.AccountID,
			&r.Description,
			&r.Created,
		); err != nil {
			return nil, err
		}

		ad, err := GetAdvertisement(adId)
		if err != nil {
			return nil, err
		}
		r.Ad = *ad

		out = append(out, r)
	}

	return out, rows.Err()
}

func FinishReport(report *utils.Report) error {
	stmt, err := utils.PrepareStmt(dat, "DELETE FROM reports WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(report.ID)
	if err != nil {
		return err
	}

	return nil
}

func init() {
	ads, err := ListAllAdvertisements()
	if err != nil {
		log.Error("Failed to initialize ads cache: %s", err.Error())
	} else {
		currentAds = &ads
		log.Info("Initialized ads cache with %d ads", len(ads))
	}
}
