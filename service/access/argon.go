package access

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"service/log"
	"service/utils"

	"github.com/patrickmn/go-cache"
)

var argonCache = cache.New(6*time.Hour, 10*time.Minute)

func UpsertArgonUser(user utils.ArgonUser) error {
	argonCache.Set(fmt.Sprintf("%d", user.Account), user, cache.DefaultExpiration)
	log.Debug("Argon cache entry added for account %d, total entries: %d", user.Account, argonCache.ItemCount())

	stmt, err := utils.PrepareStmt(utils.Db(), "INSERT INTO argon (account_id, authtoken) VALUES (?, ?) ON DUPLICATE KEY UPDATE authtoken = VALUES (authtoken), valid_at = CURRENT_TIMESTAMP")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(user.Account, user.Token)
	return err
}

func ValidateArgonUser(user utils.ArgonUser) (bool, error) {
	if val, found := argonCache.Get(fmt.Sprintf("%d", user.Account)); found {
		cUser := val.(utils.ArgonUser)

		if cUser.Account == user.Account && cUser.Token == user.Token {
			return true, nil
		}
	}

	stmt, err := utils.PrepareStmt(utils.Db(), "SELECT account_id, authtoken, valid_at FROM argon WHERE account_id = ?")
	if err != nil {
		return false, err
	} else {
		log.Debug("Prepared argon database statement for account of ID %v", user.Account)
	}

	var dbUser utils.ArgonUser
	var validAt time.Time
	if row := stmt.QueryRow(user.Account); row != nil {
		if err := row.Scan(&dbUser.Account, &dbUser.Token, &validAt); err != nil {
			log.Error(err.Error())
		} else if time.Since(validAt) < 24*time.Hour && dbUser.Account == user.Account && dbUser.Token == user.Token {
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

	log.Debug("Sending request to Argon server: %s", u.String())
	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	} else {
		log.Debug("Argon status code received: %d, for account of ID %v", resp.StatusCode, user.Account)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}
	log.Debug("Argon response body: %s", string(bodyBytes))

	var valid utils.ArgonValidation
	if err := json.Unmarshal(bodyBytes, &valid); err != nil {
		return false, err
	} else {
		log.Debug("Argon status of account of ID %v retrieved", user.Account)
	}

	if valid.Valid {
		log.Info("Argon status of account of ID %v is valid", user.Account)
		UpsertArgonUser(user)
		return true, nil
	}

	log.Error("Argon status of account of ID %v is invalid for %s", user.Account, valid.Cause)
	return false, fmt.Errorf("cause: %s", valid.Cause)
}

func DeleteArgonUser(accountID int) {
	argonCache.Delete(fmt.Sprintf("%d", accountID))
}
