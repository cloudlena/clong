package mysql

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudlena/clong/internal/app/clong"
)

// Add adds a new score to the DB.
func (s *ScoreStore) Add(ctx context.Context, scr *clong.Score) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO score
		(player_id, player_name, final_score, color)
		VALUES(?,?,?,?)
	`)
	if err != nil {
		return fmt.Errorf("error preparing DB statement: %w", err)
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Fatal(fmt.Errorf("error closing DB statement: %w", err))
		}
	}()
	_, err = stmt.Exec(scr.Player.ID, scr.Player.Name, scr.FinalScore, scr.Color)
	if err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}
