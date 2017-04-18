package main

// Message is a websocket message sent from a client
type Message struct {
	Color     string  `json:"color"`
	PosX      float64 `json:"posX"`
	PosY      float64 `json:"posY"`
	VelocityX float64 `json:"velocityX"`
	VelocityY float64 `json:"velocityY"`
}
