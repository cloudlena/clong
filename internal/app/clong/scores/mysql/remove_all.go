package mysql

import (
	"context"

	"github.com/pkg/errors"
)

// RemoveAll removes all scores from the DB.
func (db DB) RemoveAll(_ context.Context) error {
	_, err := db.session.Exec("DELETE FROM score")
	if err != nil {
		return errors.Wrap(err, "error deleting scores from DB")
	}
	return nil
}
