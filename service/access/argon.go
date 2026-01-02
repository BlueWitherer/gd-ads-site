package access

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	"service/log"
	"service/utils"

	"github.com/patrickmn/go-cache"
)

var argonCache = cache.New(6*time.Hour, 10*time.Minute)
var invalids = cache.New(1*time.Hour, 10*time.Minute)

func getToken() (string, error) {
	token := os.Getenv("ARGON")
	if token == "" {
		return "", fmt.Errorf("env for argon token is not defined!")
	} else {
		return token, nil
	}
}

func ReportBanArgonUser(report *utils.Report, banned bool) error {
	stmt, err := utils.PrepareStmt(utils.Db(), "UPDATE argon SET report_banned = ? WHERE account_id = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(banned, report.AccountID)
	return err
}

func GetArgonUser(id int) (*utils.ArgonUser, error) {
	stmt, err := utils.PrepareStmt(utils.Db(), "SELECT * FROM argon WHERE account_id = ?")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	user := new(utils.ArgonUser)
	err = stmt.QueryRow(id).Scan(&user.Account, &user.Token, &user.ReportBanned, &user.ValidAt)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func UpsertArgonUser(user *utils.ArgonUser) error {
	argonCache.Set(fmt.Sprintf("%d", user.Account), user, cache.DefaultExpiration)
	log.Debug("Argon cache entry added for account %d, total entries: %d", user.Account, argonCache.ItemCount())

	stmt, err := utils.PrepareStmt(utils.Db(), "INSERT INTO argon (account_id, authtoken) VALUES (?, ?) ON DUPLICATE KEY UPDATE authtoken = VALUES (authtoken), valid_at = CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Account, user.Token)
	return err
}

func ValidateArgonUser(user *utils.ArgonUser) (bool, error) {
	if val, found := invalids.Get(fmt.Sprintf("%d", user.Account)); found {
		return false, fmt.Errorf("Argon token %s is invalid", val.(string))
	}

	if val, found := argonCache.Get(fmt.Sprintf("%d", user.Account)); found {
		cUser := val.(*utils.ArgonUser)

		if cUser.Account == user.Account && cUser.Token == user.Token {
			return true, nil
		}
	}

	stmt, err := utils.PrepareStmt(utils.Db(), "SELECT * FROM argon WHERE account_id = ?")
	if err != nil {
		return false, err
	} else {
		log.Debug("Prepared argon database statement for account of ID %v", user.Account)
	}
	defer stmt.Close()

	dbUser := new(utils.ArgonUser)
	if row := stmt.QueryRow(user.Account); row != nil {
		if err := row.Scan(&dbUser.Account, &dbUser.Token, &dbUser.ReportBanned, &dbUser.ValidAt); err != nil {
			log.Error(err.Error())
		} else if time.Since(dbUser.ValidAt) < 24*time.Hour && dbUser.Account == user.Account && dbUser.Token == user.Token {
			return true, nil
		}
	}

	u, err := url.Parse("https://argon.globed.dev/v1/validation/check")
	if err != nil {
		return false, err
	} else {
		log.Debug("Argon URL parsed for account of ID %v", user.Account)
	}

	q := u.Query()
	q.Set("account_id", fmt.Sprintf("%v", user.Account))
	q.Set("authtoken", user.Token)
	u.RawQuery = q.Encode()

	log.Debug("Argon validation parameters: account_id=%d (type check: %T), authtoken length=%d", user.Account, user.Account, len(user.Token))
	log.Debug("Full Argon URL being requested: %s", u.String())

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return false, err
	} else {
		log.Debug("Argon request object constructed for account of ID %v", user.Account)
	}

	req.Header.Set("User-Agent", "PlayerAdvertisements/1.0")

	argon, err := getToken()
	if err != nil {
		log.Warn("Failed to get Argon API token: %s", err.Error())
	} else {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", argon))
	}

	log.Debug("Sending request to Argon server: %s", u.String())
	client := &http.Client{Timeout: 15 * time.Second}

	resp, reqErr := client.Do(req)
	if reqErr != nil {
		return false, reqErr
	}
	defer resp.Body.Close()

	log.Debug("Argon status code received: %d, for account of ID %v", resp.StatusCode, user.Account)

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	log.Debug("Argon response body: %s", string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return false, fmt.Errorf("argon server returned status code %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var valid utils.ArgonValidation
	if err := json.Unmarshal(bodyBytes, &valid); err != nil {
		return false, fmt.Errorf("failed to parse argon response: %v, body: %s", err, string(bodyBytes))
	} else {
		log.Debug("Argon status of account of ID %v retrieved", user.Account)
	}

	if valid.Valid {
		log.Info("Argon status of account of ID %v is valid", user.Account)
		UpsertArgonUser(user)

		return true, nil
	}

	invalids.Set(fmt.Sprintf("%d", user.Account), user.Token, cache.DefaultExpiration)
	return false, fmt.Errorf("cause: %s", valid.Cause)
}
