package clong

// Control is a WebSocket message sent from a controller to control a screen.
type Control struct {
	Player    string  `json:"player"`
	Color     string  `json:"color"`
	PosX      float64 `json:"posX"`
	PosY      float64 `json:"posY"`
	VelocityX float64 `json:"velocityX"`
	VelocityY float64 `json:"velocityY"`
}
