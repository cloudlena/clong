package clong

import "github.com/mastertinner/clong/internal/app/clong/users"

// event is a WebSocket message sent from a screen to a controller.
type event struct {
	Type   string     `json:"type"`
	Player users.User `json:"player"`
	Points int64      `json:"points"`
}
