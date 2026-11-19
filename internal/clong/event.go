package clong

// EventBallDone is sent when a ball hit a target or left the screen.
const EventBallDone = "BALL_DONE"

// Event is a message sent to a controller.
type Event struct {
	Type   string `json:"type"`
	Player User   `json:"player"`
	Points int64  `json:"points"`
}
