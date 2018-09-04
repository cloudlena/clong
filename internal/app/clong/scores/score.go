package scores

import (
	"context"

	"github.com/mastertinner/clong/internal/app/clong/users"
)

// Score is the score of a player.
type Score struct {
	ID         string     `json:"id"`
	Player     users.User `json:"player"`
	FinalScore int64      `json:"finalScore"`
	Color      string     `json:"color"`
}

// ScoreStore is a store of scores.
type ScoreStore interface {
	Scores(ctx context.Context) ([]Score, error)
	Add(ctx context.Context, data Score) error
	RemoveAll(ctx context.Context) error
}
