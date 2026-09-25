package pg

import (
	"context"
	"fmt"

	"github.com/cloudlena/clong/internal/clong"
)

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
