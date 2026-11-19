package clong

// Control types handled by the service.
const (
	ControlGameStarted  = "GAME_STARTED"
	ControlGameFinished = "GAME_FINISHED"
)

// Control is a message sent from a controller.
type Control struct {
	Type      string  `json:"type"`
	Player    User    `json:"player"`
	Color     string  `json:"color"`
	PosX      float64 `json:"posX"`
	VelocityX float64 `json:"velocityX"`
	VelocityY float64 `json:"velocityY"`
}
