package database

import (
	"fmt"
	"os"
	"path/filepath"

	"service/utils"
)

func GetUser(id string) (utils.User, error) {
	if id == "" {
		return utils.User{}, fmt.Errorf("empty user id")
	}

	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM users WHERE id = ?")
	if err != nil {
		return utils.User{}, err
	}

	var user utils.User
	err = stmt.QueryRow(id).Scan(
		&user.ID,
		&user.Username,
		&user.AvatarURL,
		&user.TotalViews,
		&user.TotalClicks,
		&user.IsAdmin,
		&user.IsStaff,
		&user.Verified,
		&user.Banned,
		&user.BoostCount,
		&user.Created,
		&user.Updated,
	)
	if err != nil {
		return utils.User{}, err
	}

	return user, nil
}

func GetAllUsers() ([]utils.User, error) {
	stmt, err := utils.PrepareStmt(dat, "SELECT * FROM users ORDER BY id DESC")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	users, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer users.Close()

	var out []utils.User
	for users.Next() {
		var u utils.User
		if err := users.Scan(
			&u.ID,
			&u.Username,
			&u.AvatarURL,
			&u.TotalClicks,
			&u.TotalViews,
			&u.IsAdmin,
			&u.IsStaff,
			&u.Verified,
			&u.Banned,
			&u.BoostCount,
			&u.Created,
			&u.Updated,
		); err != nil {
			return nil, err
		}

		out = append(out, u)
	}

	return out, users.Err()
}

// inserts a new user or updates username if it already exists.
func UpsertUser(id string, username string, avatarUrl string) error {
	if id == "" {
		return fmt.Errorf("empty user id")
	}

	stmt, err := utils.PrepareStmt(dat, "INSERT INTO users (id, username, avatar_url) VALUES (?, ?, ?) ON DUPLICATE KEY UPDATE username = VALUES (username), avatar_url = VALUES (avatar_url), updated_at = CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id, username, avatarUrl)
	return err
}

// increments total_views or total_clicks for an ad
func IncrementAdStat(adID int64, statType utils.AdEvent) error {
	query := fmt.Sprintf("UPDATE advertisements SET %s = %s + 1 WHERE ad_id = ?", statType, statType)

	stmt, err := utils.PrepareStmt(dat, query)
	if err != nil {
		return fmt.Errorf("failed to prepare increment query: %w", err)
	}
	defer stmt.Close()

	if _, err := stmt.Exec(adID); err != nil {
		return fmt.Errorf("failed to increment %s for ad %d: %w", statType, adID, err)
	}

	return nil
}

// increments total_views or total_clicks for a user
func IncrementUserStats(userId string, viewsDelta int, clicksDelta int) error {
	if userId == "" {
		return fmt.Errorf("empty user id")
	}

	stmt, err := utils.PrepareStmt(dat, "UPDATE users SET total_views = total_views + ?, total_clicks = total_clicks + ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(viewsDelta, clicksDelta, userId)
	return err
}

func BanUser(id string) (utils.User, error) {
	// delete all advertisements associated with the user
	deleteAdsStmt, err := utils.PrepareStmt(dat, "SELECT * FROM advertisements WHERE user_id = ?")
	if err != nil {
		return utils.User{}, err
	}
	defer deleteAdsStmt.Close()

	rows, err := deleteAdsStmt.Query(id)
	if err != nil {
		return utils.User{}, err
	}

	ads := make([]utils.Ad, 0)
	for rows.Next() {
		var r utils.Ad
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
			return utils.User{}, err
		}

		ads = append(ads, r)
	}

	user, err := GetUser(id)
	if err != nil {
		return utils.User{}, err
	}

	for _, a := range ads {
		t, err := utils.AdTypeFromInt(a.Type)
		if err != nil {
			return utils.User{}, err
		}

		adDir := filepath.Join("..", "ad_storage", string(t), fmt.Sprintf("%s-%d.webp", a.UserID, a.AdID))
		err = os.Remove(adDir)
		if err != nil {
			return utils.User{}, err
		}
	}

	// ban the user
	stmt, err := utils.PrepareStmt(dat, "UPDATE users SET banned = TRUE WHERE id = ?")
	if err != nil {
		return utils.User{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return utils.User{}, err
	}

	return user, nil
}

func UnbanUser(id string) (utils.User, error) {
	// unban the user
	stmt, err := utils.PrepareStmt(dat, "UPDATE users SET banned = FALSE WHERE id = ?")
	if err != nil {
		return utils.User{}, err
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return utils.User{}, err
	}

	return GetUser(id)
}

func UserLeaderboard(stat utils.StatBy, page uint64, maxPerPage uint64) ([]utils.User, error) {
	stmt, err := utils.PrepareStmt(dat, fmt.Sprintf("SELECT * FROM users WHERE banned = FALSE ORDER BY %s DESC", stat))
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var out []utils.User
	for rows.Next() {
		var r utils.User

		if err := rows.Scan(
			&r.ID,
			&r.Username,
			&r.AvatarURL,
			&r.TotalViews,
			&r.TotalClicks,
			&r.IsAdmin,
			&r.IsStaff,
			&r.Verified,
			&r.Banned,
			&r.BoostCount,
			&r.Created,
			&r.Updated,
		); err != nil {
			return nil, err
		}

		out = append(out, r)
	}

	start := page * maxPerPage
	end := start + maxPerPage

	if start >= uint64(len(out)) {
		return []utils.User{}, nil
	}

	if end > uint64(len(out)) {
		end = uint64(len(out))
	}

	return out[start:end], nil
}
