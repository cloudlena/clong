package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Hub is a WebSocket messaging hub
type Hub struct {
	Clients          map[*websocket.Conn]bool
	RegisterClient   chan *websocket.Conn
	UnregisterClient chan *websocket.Conn
	Screens          map[*websocket.Conn]bool
	RegisterScreen   chan *websocket.Conn
	UnregisterScreen chan *websocket.Conn
	Broadcast        chan Message
}

// newHub creates a new messaging hub
func newHub() *Hub {
	return &Hub{
		Clients:          make(map[*websocket.Conn]bool),
		RegisterClient:   make(chan *websocket.Conn),
		UnregisterClient: make(chan *websocket.Conn),
		Screens:          make(map[*websocket.Conn]bool),
		RegisterScreen:   make(chan *websocket.Conn),
		UnregisterScreen: make(chan *websocket.Conn),
		Broadcast:        make(chan Message),
	}
}

// run runs the messaging hub in a forever loop
func (h *Hub) run() {
	for {
		select {
		case c := <-h.RegisterClient:
			h.Clients[c] = true
			log.Printf("client added. (%v connected)", len(h.Clients))

		case c := <-h.UnregisterClient:
			delete(h.Clients, c)
			log.Printf("client removed (%v connected)", len(h.Clients))

		case s := <-h.RegisterScreen:
			h.Screens[s] = true
			log.Printf("screen added. (%v connected)", len(h.Screens))

		case s := <-h.UnregisterScreen:
			delete(h.Screens, s)
			log.Printf("screen removed (%v connected)", len(h.Screens))

		case msg := <-h.Broadcast:
			for s := range h.Screens {
				err := s.WriteJSON(msg)
				if err != nil {
					s.Close()
					h.UnregisterScreen <- s
				}
			}
		}
	}
}
