package mysql

import (
	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// CreateScore creates a new score.
func (db DB) CreateScore(in clong.Score) error {
	smtp, err := db.Prepare("INSERT INTO scores(playerID, playerName, finalScore, color) VALUES(?,?,?,?)")
	if err != nil {
		return errors.Wrap(err, "error preparing create score DB statement")
	}
	_, err = smtp.Exec(in.Player.ID, in.Player.Name, in.FinalScore, in.Color)
	if err != nil {
		return errors.Wrap(err, "error executing create score DB statement")
	}
	return nil
}
