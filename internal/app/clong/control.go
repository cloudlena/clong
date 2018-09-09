package clong

// Control is a message sent from a controller.
type Control struct {
	Type       string  `json:"type"`
	Player     User    `json:"player"`
	Color      string  `json:"color"`
	PosX       float64 `json:"posX"`
	PosY       float64 `json:"posY"`
	VelocityX  float64 `json:"velocityX"`
	VelocityY  float64 `json:"velocityY"`
	FinalScore int64   `json:"finalScore"`
}
