package clong

// Event is a message sent to a controller.
type Event struct {
	Type   string `json:"type"`
	Player User   `json:"player"`
	Points int64  `json:"points"`
}
