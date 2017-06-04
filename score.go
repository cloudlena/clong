package clong

// Score is the score of a player.
type Score struct {
	ID         string `json:"id"`
	Player     User   `json:"player"`
	FinalScore int    `json:"finalScore"`
	Color      string `json:"color"`
}
