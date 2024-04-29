package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudlena/clong/internal/clong"
)

// ListAll retrieves all scores from the DB.
func (s *ScoreStore) ListAll(ctx context.Context) ([]*clong.Score, error) {
	rows, err := s.db.QueryContext(ctx, "SELECT score_id, player_id, player_name, final_score, color FROM score")
	if err != nil {
		return nil, fmt.Errorf("error querying DB: %w", err)
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("error closing DB rows: %w", err))
		}
	}()

	scrs := make([]*clong.Score, 0)
	for rows.Next() {
		var s clong.Score
		err = rows.Scan(&s.ID, &s.Player.ID, &s.Player.Name, &s.FinalScore, &s.Color)
		if err != nil {
			return nil, fmt.Errorf("error scanning DB rows: %w", err)
		}
		scrs = append(scrs, &s)
	}
	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error in DB rows: %w", err)
	}

	return scrs, nil
}
