package clong

// Event is a WebSocket message sent from a screen to a controller.
type Event struct {
	MsgType string `json:"msgType"`
	Player  User   `json:"player"`
	Points  int64  `json:"points"`
}
