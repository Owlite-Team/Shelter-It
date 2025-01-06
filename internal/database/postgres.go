package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	"time"
)

type Database struct {
	DB *sql.DB
}

func NewDB(connectionString string) (*Database, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Config conn pool
	db.SetMaxOpenConns(25)                 // Simultaneous conn
	db.SetMaxIdleConns(5)                  // Keep conn ready
	db.SetConnMaxLifetime(5 * time.Minute) // Refresh conn periodically

	// Verify conn is working
	if err = db.Ping(); err != nil {
		return nil, err
	}

	return &Database{db}, nil
}
