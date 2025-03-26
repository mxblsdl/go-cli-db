package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq" // PostgreSQL driver
)

var db *sql.DB

// Connect establishes a connection to the PostgreSQL database.
func Connect(connStr string) {
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatalf("Error pinging the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")
}

// Query executes a SQL query and returns the result.
func Query(query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

