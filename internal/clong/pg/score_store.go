// Package pg stores data in a PostgreSQL database.
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

// NewScoreStore creates a new score store, creating its table if needed.
func NewScoreStore(ctx context.Context, db *sql.DB) (*ScoreStore, error) {
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("error pinging DB: %w", err)
	}

	if _, err := db.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS score (
		score_id SERIAL NOT NULL PRIMARY KEY,
		player_id VARCHAR(36) NOT NULL,
		player_name VARCHAR(30) NOT NULL,
		final_score INT NOT NULL,
		color VARCHAR(7) NOT NULL
	)`); err != nil {
		return nil, fmt.Errorf("error executing DB statement: %w", err)
	}

	return &ScoreStore{db: db}, nil
}

// ListTop retrieves the best scores from the DB, best first. Ties go to the earlier score.
func (s *ScoreStore) ListTop(ctx context.Context, limit int) ([]clong.Score, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT player_name, final_score, color FROM score
		ORDER BY final_score DESC, score_id
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("error querying DB: %w", err)
	}
	defer func() {
		if err := rows.Close(); err != nil {
			log.Printf("error closing DB rows: %v", err)
		}
	}()

	// Not nil, so it's encoded as an empty JSON array
	scrs := make([]clong.Score, 0)
	for rows.Next() {
		var scr clong.Score
		if err := rows.Scan(&scr.PlayerName, &scr.FinalScore, &scr.Color); err != nil {
			return nil, fmt.Errorf("error scanning DB rows: %w", err)
		}
		scrs = append(scrs, scr)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error in DB rows: %w", err)
	}

	return scrs, nil
}

// Add adds a new score to the DB.
func (s *ScoreStore) Add(ctx context.Context, scr clong.Score) error {
	if _, err := s.db.ExecContext(ctx, `
		INSERT INTO score
		(player_id, player_name, final_score, color)
		VALUES ($1, $2, $3, $4)
	`, scr.PlayerID, scr.PlayerName, scr.FinalScore, scr.Color); err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}

// RemoveAll removes all scores from the DB.
func (s *ScoreStore) RemoveAll(ctx context.Context) error {
	if _, err := s.db.ExecContext(ctx, "DELETE FROM score"); err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}
