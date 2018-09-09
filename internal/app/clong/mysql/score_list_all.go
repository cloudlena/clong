package mysql

import (
	"context"
	"log"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// ListAll retrieves all scores from the DB.
func (s *scoreStore) ListAll(ctx context.Context) (scrs []*clong.Score, err error) {
	rows, err := s.db.QueryContext(ctx, "SELECT score_id, player_id, player_name, final_score, color FROM score")
	if err != nil {
		return nil, errors.Wrap(err, "error querying DB")
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB rows"))
		}
	}()
	for rows.Next() {
		var s clong.Score
		err = rows.Scan(&s.ID, &s.Player.ID, &s.Player.Name, &s.FinalScore, &s.Color)
		if err != nil {
			return nil, errors.Wrap(err, "error scanning DB rows")
		}
		scrs = append(scrs, &s)
	}
	err = rows.Err()
	if err != nil {
		return nil, errors.Wrap(err, "error in DB rows")
	}
	return scrs, nil
}
