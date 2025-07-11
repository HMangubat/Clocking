package config

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

// InitDB initializes and returns a DB connection.
func InitDB() *sql.DB {
	connStr := "postgres://postgres:123@localhost:5432/hvm?sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Fatal: could not connect to DB: %v", err)
	}
	return db
}
