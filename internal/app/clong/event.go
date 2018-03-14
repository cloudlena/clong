package clong

// Event is a WebSocket message sent from a screen to a controller.
type Event struct {
	Type   string `json:"type"`
	Player User   `json:"player"`
	Points int64  `json:"points"`
}
