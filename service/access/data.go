package access

import (
	"database/sql"
	"os"
	"time"

	"service/log"

	_ "github.com/go-sql-driver/mysql"
)

var data *sql.DB

type AdEvent string

const (
	AdEventView  AdEvent = "view"  // For views
	AdEventClick AdEvent = "click" // For clicks
)

type StatBy string

const (
	StatByAd   StatBy = "ad_id"   // Filter stats by ad
	StatByUser StatBy = "user_id" // Filter stats by user
)

// Register a new client event for an ad
func NewStat(event AdEvent, ad int64, user int64) error {
	log.Debug("Registering new " + event)
	stmt, err := data.Prepare("INSERT INTO ad_views (ad_id, user_id, timestamp) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(ad, user, time.Now())
	return err
}

func init() {
	var err error

	data, err = sql.Open("mysql", os.Getenv("DB_URI"))
	if err != nil {
		log.Error(err.Error())
		return
	}

	err = data.Ping()
	if err != nil {
		log.Error(err.Error())
		return
	}

	log.Print("MariaDB connection established.")
}
