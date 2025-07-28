package pg

import (
	"context"
	"fmt"
)

// RemoveAll removes all scores from the DB.
func (s *ScoreStore) RemoveAll(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, "DELETE FROM score")
	if err != nil {
		return fmt.Errorf("error executing DB statement: %w", err)
	}
	return nil
}
