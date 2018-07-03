// Package mysql allows to interact with a MySQL database.
package mysql

import (
	"database/sql"

	"github.com/pkg/errors"
)

// DB is a database.
type DB struct {
	session *sql.DB
}

// New creates a new database session.
func New(connString string) (DB, error) {
	sess, err := sql.Open("mysql", connString)
	if err != nil {
		return DB{}, errors.Wrap(err, "error opening DB session")
	}

	// Check if DB connection is healthy
	err = sess.Ping()
	if err != nil {
		return DB{}, errors.Wrap(err, "error pinging DB")
	}

	// Create scores table if it doesn't exist yet
	_, err = sess.Exec(`CREATE TABLE IF NOT EXISTS scores (
		id int NOT NULL AUTO_INCREMENT,
		playerID varchar(36),
		playerName varchar(30),
		finalScore int,
		created date,
		color varchar(7),
		PRIMARY KEY (id)
	)`)
	if err != nil {
		return DB{}, errors.Wrap(err, "error creating table in DB")
	}

	return DB{session: sess}, nil
}

// Close closes the underlying database session.
func (db DB) Close() error {
	return db.session.Close()
}
