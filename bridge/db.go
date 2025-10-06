package main

import (
	"database/sql"
	"fmt"
	"os"

	"bridge/log"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func InitDB() error {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		// default to user:password@tcp(host:port)/dbname
		// The user should set DATABASE_DSN to something like: user:pass@tcp(127.0.0.1:3306)/gdads
		dsn = "root:@tcp(127.0.0.1:3306)/gdads?parseTime=true"
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		return fmt.Errorf("sql.Open: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("db.Ping: %w", err)
	}

	log.Info("Connected to database")

	if err := ensureSchema(); err != nil {
		return err
	}

	return nil
}

// ensureSchema creates the users table if it doesn't exist
func ensureSchema() error {
	q := `CREATE TABLE IF NOT EXISTS users (
        id VARCHAR(64) PRIMARY KEY,
        username VARCHAR(255) NOT NULL,
        discriminator VARCHAR(16),
        avatar VARCHAR(255),
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;`

	_, err := db.Exec(q)
	if err != nil {
		return fmt.Errorf("ensureSchema: %w", err)
	}
	return nil
}

// UpsertUser inserts or updates a user record
func UpsertUser(u User) error {
	if db == nil {
		return fmt.Errorf("database not initialized")
	}

	// Use INSERT ... ON DUPLICATE KEY UPDATE for upsert
	q := `INSERT INTO users (id, username, discriminator, avatar) VALUES (?, ?, ?, ?)
    ON DUPLICATE KEY UPDATE username=VALUES(username), discriminator=VALUES(discriminator), avatar=VALUES(avatar)`

	_, err := db.Exec(q, u.ID, u.Username, u.Discriminator, u.Avatar)
	if err != nil {
		return fmt.Errorf("UpsertUser: %w", err)
	}
	return nil
}
