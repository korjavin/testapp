package repository

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

// DB is the database connection
var DB *sql.DB

// InitDB initializes the database connection
func InitDB(dsn string) error {
	var err error
	DB, err = sql.Open("sqlite", dsn)
	if err != nil {
		return err
	}
	return DB.Ping()
}

// Close closes the database connection
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
