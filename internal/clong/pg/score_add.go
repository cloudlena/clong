package pg

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudlena/clong/internal/clong"
)

// Add adds a new score to the DB.
func (s *ScoreStore) Add(ctx context.Context, scr *clong.Score) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO score
		(player_id, player_name, final_score, color)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return fmt.Errorf("error preparing DB statement: %w", err)
	}
	defer func() {
		if cErr := stmt.Close(); cErr != nil {
			log.Printf("error closing DB statemet :%v\n", cErr)
		}
	}()

	_, err = stmt.ExecContext(ctx, scr.Player.ID, scr.Player.Name, scr.FinalScore, scr.Color)
	if err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}

	return nil
}
