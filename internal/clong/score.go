package clong

import "context"

// Score is a score of a player.
type Score struct {
	ID         string `json:"id"`
	Player     User   `json:"player"`
	FinalScore int64  `json:"finalScore"`
	Color      string `json:"color"`
}

// ScoreStore is a store of scores.
type ScoreStore interface {
	ListAll(ctx context.Context) ([]*Score, error)
	Add(ctx context.Context, scr *Score) error
	RemoveAll(ctx context.Context) error
}
