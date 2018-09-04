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

// Make creates a new database session.
func Make(session *sql.DB) (DB, error) {
	// Check if DB connection is healthy
	err := session.Ping()
	if err != nil {
		return DB{}, errors.Wrap(err, "error pinging DB")
	}

	// Create score table if it doesn't exist yet
	_, err = session.Exec(`CREATE TABLE IF NOT EXISTS score (
		score_id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		player_id VARCHAR(36) NOT NULL,
		player_name VARCHAR(30) NOT NULL,
		final_score INT(11) UNSIGNED NOT NULL,
		color VARCHAR(7) NOT NULL,
		PRIMARY KEY (score_id)
	)`)
	if err != nil {
		return DB{}, errors.Wrap(err, "error creating table in DB")
	}

	return DB{session: session}, nil
}

// Close closes the underlying database session.
func (db DB) Close() error {
	err := db.session.Close()
	if err != nil {
		return errors.Wrap(err, "error closing DB connection")
	}
	return nil
}
