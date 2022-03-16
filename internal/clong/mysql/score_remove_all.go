package mysql

import (
	"context"
	"fmt"
)

// RemoveAll removes all scores from the DB.
func (s *ScoreStore) RemoveAll(_ context.Context) error {
	_, err := s.db.Exec("DELETE FROM score")
	if err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}
