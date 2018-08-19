package mysql

import (
	"context"

	"github.com/pkg/errors"
)

// DeleteScores deletes all scores from the DB.
func (db DB) DeleteScores(_ context.Context) error {
	_, err := db.session.Exec("DELETE FROM score")
	if err != nil {
		return errors.Wrap(err, "error deleting scores from DB")
	}
	return nil
}
