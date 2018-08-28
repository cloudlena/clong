package clong

import "github.com/mastertinner/clong/internal/app/clong/users"

// control is a WebSocket message sent from a controller to control a screen.
type control struct {
	Type       string     `json:"type"`
	Player     users.User `json:"player"`
	Color      string     `json:"color"`
	PosX       float64    `json:"posX"`
	PosY       float64    `json:"posY"`
	VelocityX  float64    `json:"velocityX"`
	VelocityY  float64    `json:"velocityY"`
	FinalScore int64      `json:"finalScore"`
}
