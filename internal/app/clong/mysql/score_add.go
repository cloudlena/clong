package mysql

import (
	"context"
	"log"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// Add adds a new score to the DB.
func (s *scoreStore) Add(ctx context.Context, scr *clong.Score) error {
	stmt, err := s.db.PrepareContext(ctx, `
		INSERT INTO score
		(player_id, player_name, final_score, color)
		VALUES(?,?,?,?)
	`)
	if err != nil {
		return errors.Wrap(err, "error preparing DB statement")
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB statement"))
		}
	}()
	_, err = stmt.Exec(scr.Player.ID, scr.Player.Name, scr.FinalScore, scr.Color)
	if err != nil {
		return errors.Wrap(err, "error executing DB statement")
	}
	return nil
}
