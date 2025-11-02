package utils

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"service/log"

	_ "github.com/go-sql-driver/mysql"
)

// Concurrent database connection
var data *sql.DB

// safely prepare the sql statement
func PrepareStmt(db *sql.DB, sql string) (*sql.Stmt, error) {
	if db != nil {
		log.Debug("Preparing connection for statement %s", sql)
		return db.Prepare(sql)
	} else {
		return nil, fmt.Errorf("database connection non-existent")
	}
}

func Db() *sql.DB {
	return data
}

// initializeSchema reads and executes the schema.sql file to create tables if they don't exist
func initializeSchema() error {
	schemaPath := filepath.Join("..", "database", "schema.sql")
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

func init() {
	var err error

	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
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
