package clong

import "context"

// Score is a score of a player.
type Score struct {
	PlayerID   string `json:"-"` // not public, since it would allow others to play as the player
	PlayerName string `json:"playerName"`
	FinalScore int64  `json:"finalScore"`
	Color      string `json:"color"`
}

// ScoreStore is a store of scores.
type ScoreStore interface {
	// ListTop returns the best scores, best first.
	ListTop(ctx context.Context, limit int) ([]Score, error)
	Add(ctx context.Context, scr Score) error
	RemoveAll(ctx context.Context) error
}
