package mysql

import (
	"context"

	"github.com/pkg/errors"
)

// RemoveAll removes all scores from the DB.
func (s *scoreStore) RemoveAll(_ context.Context) error {
	_, err := s.db.Exec("DELETE FROM score")
	if err != nil {
		return errors.Wrap(err, "error executing DB statement")
	}
	return nil
}
