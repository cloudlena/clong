package clong

// event is a WebSocket message sent from a screen to a controller.
type event struct {
	Type   string `json:"type"`
	Player user   `json:"player"`
	Points int64  `json:"points"`
}
