// Package pg allows to interact with a PostgreSQL database.
package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/cloudlena/clong/internal/clong"
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
		score_id SERIAL NOT NULL PRIMARY KEY,
		player_id VARCHAR(36) NOT NULL,
		player_name VARCHAR(30) NOT NULL,
		final_score INT NOT NULL,
		color VARCHAR(7) NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("error executing DB statement: %w", err)
	}

	return &ScoreStore{db: db}, nil
}

// ListAll retrieves all scores from the DB.
func (s *ScoreStore) ListAll(ctx context.Context) ([]*clong.Score, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT score_id, player_id, player_name, final_score, color FROM score")
	if err != nil {
		return nil, fmt.Errorf("error querying DB: %w", err)
	}
	defer func() {
		if cErr := rows.Close(); cErr != nil {
			log.Printf("error closing DB rows: %v\n", cErr)
		}
	}()

	scrs := make([]*clong.Score, 0)
	for rows.Next() {
		var scr clong.Score
		err = rows.Scan(&scr.ID, &scr.Player.ID, &scr.Player.Name, &scr.FinalScore, &scr.Color)
		if err != nil {
			return nil, fmt.Errorf("error scanning DB rows: %w", err)
		}
		scrs = append(scrs, &scr)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error in DB rows: %w", err)
	}

	return scrs, nil
}

// Add adds a new score to the DB.
func (s *ScoreStore) Add(ctx context.Context, scr *clong.Score) error {
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO score
		(player_id, player_name, final_score, color)
		VALUES ($1, $2, $3, $4)
	`, scr.Player.ID, scr.Player.Name, scr.FinalScore, scr.Color)
	if err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}

// RemoveAll removes all scores from the DB.
func (s *ScoreStore) RemoveAll(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM score")
	if err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}
