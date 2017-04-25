package main

import (
	"log"

	"github.com/gorilla/websocket"
)

// Hub is a WebSocket messaging hub
type Hub struct {
	Controllers          map[*websocket.Conn]bool
	Screens              map[*websocket.Conn]bool
	RegisterController   chan *websocket.Conn
	UnregisterController chan *websocket.Conn
	RegisterScreen       chan *websocket.Conn
	UnregisterScreen     chan *websocket.Conn
	Controls             chan Control
	Events               chan Event
}

// newHub creates a new messaging hub
func newHub() *Hub {
	return &Hub{
		Controllers:          make(map[*websocket.Conn]bool),
		Screens:              make(map[*websocket.Conn]bool),
		RegisterController:   make(chan *websocket.Conn),
		UnregisterController: make(chan *websocket.Conn),
		RegisterScreen:       make(chan *websocket.Conn),
		UnregisterScreen:     make(chan *websocket.Conn),
		Controls:             make(chan Control),
		Events:               make(chan Event),
	}
}

// run runs the messaging hub in a forever loop
func (h *Hub) run() {
	for {
		select {
		case c := <-h.RegisterController:
			h.Controllers[c] = true
			log.Printf("controller added (%v connected)", len(h.Controllers))

		case c := <-h.UnregisterController:
			delete(h.Controllers, c)
			log.Printf("controller removed (%v connected)", len(h.Controllers))

		case s := <-h.RegisterScreen:
			h.Screens[s] = true
			log.Printf("screen added (%v connected)", len(h.Screens))

		case s := <-h.UnregisterScreen:
			delete(h.Screens, s)
			log.Printf("screen removed (%v connected)", len(h.Screens))

		case c := <-h.Controls:
			for s := range h.Screens {
				err := s.WriteJSON(c)
				if err != nil {
					err = s.Close()
					if err != nil {
						log.Fatal(err)
					}
					h.UnregisterScreen <- s
				}
			}

		case e := <-h.Events:
			for c := range h.Controllers {
				err := c.WriteJSON(e)
				if err != nil {
					err = c.Close()
					if err != nil {
						log.Fatal(err)
					}
					h.UnregisterController <- c
				}
			}
		}
	}
}