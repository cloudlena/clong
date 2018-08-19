package mysql

import (
	"context"
	"log"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// CreateScore creates a new score.
func (db DB) CreateScore(ctx context.Context, in clong.Score) error {
	stmt, err := db.session.PrepareContext(ctx, "INSERT INTO score (player_id, player_name, final_score, color) VALUES(?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, "error preparing create score DB statement")
	}
	defer func() {
		err := stmt.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB statement"))
		}
	}()
	_, err = stmt.Exec(in.Player.ID, in.Player.Name, in.FinalScore, in.Color)
	if err != nil {
		return errors.Wrap(err, "error executing create score DB statement")
	}
	return nil
}
