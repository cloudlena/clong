package clong

// Score is the score of a player.
type Score struct {
	ID         string `json:"id"`
	Player     user   `json:"player"`
	FinalScore int    `json:"finalScore"`
	Color      string `json:"color"`
}

// ScoreStore is a store of scores.
type ScoreStore interface {
	Scores() ([]Score, error)
	CreateScore(Score) error
	DeleteScores() error
}
