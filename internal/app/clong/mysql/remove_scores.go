package mysql

import (
	"log"

	"github.com/pkg/errors"
)

// RemoveScores removes all scores from the DB.
func (db DB) RemoveScores() error {
	rows, err := db.Query("DELETE FROM scores")
	if err != nil {
		return errors.Wrap(err, "error removing scores from DB")
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB rows"))
		}
	}()

	return nil
}
