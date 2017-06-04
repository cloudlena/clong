package clong

import (
	"database/sql"

	"github.com/pkg/errors"
)

// DB is a database.
type DB struct {
	*sql.DB
}

// NewDB creates a new database session.
func NewDB(connString string) (DB, error) {
	db, err := sql.Open("mysql", connString)
	if err != nil {
		return DB{}, errors.Wrap(err, "error opening DB session")
	}

	// Check if DB connection is healthy
	err = db.Ping()
	if err != nil {
		return DB{}, errors.Wrap(err, "error pinging DB")
	}

	// Create scores table if it doesn't exist yet
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scores (id int NOT NULL AUTO_INCREMENT, playerID varchar(36), playerName varchar(15), finalScore int, created date, color varchar(7), PRIMARY KEY (id))")
	if err != nil {
		return DB{}, errors.Wrap(err, "error creating scores table")
	}

	return DB{db}, nil
}
