package clong

// Event is a WebSocket message sent from a screen to control a controller.
type Event struct {
	Type   string `json:"type"`
	Player string `json:"player"`
	Points int64  `json:"points"`
}
