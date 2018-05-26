// Package mysql allows to interact with a MySQL storage backend.
package mysql // import "github.com/mastertinner/clong/internal/app/clong/mysql"

import (
	"database/sql"

	"github.com/pkg/errors"
)

// DB is a database.
type DB struct {
	*sql.DB
}

// New creates a new database session.
func New(connString string) (DB, error) {
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
	_, err = db.Exec("CREATE TABLE IF NOT EXISTS scores (id int NOT NULL AUTO_INCREMENT, playerID varchar(36), playerName varchar(30), finalScore int, created date, color varchar(7), PRIMARY KEY (id))")
	if err != nil {
		return DB{}, errors.Wrap(err, "error creating table in DB")
	}

	return DB{db}, nil
}
