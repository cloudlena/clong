// Package mysql allows to interact with a MySQL database.
package mysql

import (
	"database/sql"
	"fmt"
)

// ScoreStore is a score store.
type ScoreStore struct {
	db *sql.DB
}

// NewScoreStore creates a new score store.
func NewScoreStore(db *sql.DB) (*ScoreStore, error) {
	// Check if DB connection is healthy
	err := db.Ping()
	if err != nil {
		return nil, fmt.Errorf("error pinging DB: %w", err)
	}

	// Create score table if it doesn't exist yet
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS score (
		score_id INT(11) UNSIGNED NOT NULL AUTO_INCREMENT,
		player_id VARCHAR(36) NOT NULL,
		player_name VARCHAR(30) NOT NULL,
		final_score INT(11) UNSIGNED NOT NULL,
		color VARCHAR(7) NOT NULL,
		PRIMARY KEY (score_id)
	)`)
	if err != nil {
		return nil, fmt.Errorf("error executing DB statement: %w", err)
	}

	return &ScoreStore{db: db}, nil
}
