package mysql

import (
	"log"

	"github.com/pkg/errors"
)

// DeleteScores deletes all scores from the DB.
func (db DB) DeleteScores() error {
	rows, err := db.session.Query("DELETE FROM score")
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
