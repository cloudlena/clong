package clong

import (
	"log"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// Hub is a WebSocket messaging hub.
type Hub struct {
	RegisterController   chan *websocket.Conn
	UnregisterController chan *websocket.Conn
	RegisterScreen       chan *websocket.Conn
	UnregisterScreen     chan *websocket.Conn
	Controls             chan Control
	Events               chan Event
	controllers          map[*websocket.Conn]bool
	screens              map[*websocket.Conn]bool
	db                   DB
}

// NewHub creates a new messaging hub.
func NewHub(db DB) *Hub {
	return &Hub{
		RegisterController:   make(chan *websocket.Conn),
		UnregisterController: make(chan *websocket.Conn),
		RegisterScreen:       make(chan *websocket.Conn),
		UnregisterScreen:     make(chan *websocket.Conn),
		Controls:             make(chan Control),
		Events:               make(chan Event),
		controllers:          make(map[*websocket.Conn]bool),
		screens:              make(map[*websocket.Conn]bool),
		db:                   db,
	}
}

// Run runs the messaging hub in a forever loop as a goruotine.
func (h *Hub) Run() { // nolint: gocyclo
	go func() {
		for {
			select {
			case c := <-h.RegisterController:
				h.controllers[c] = true
				log.Printf("controller registered (%v connected)", len(h.controllers))

			case c := <-h.UnregisterController:
				delete(h.controllers, c)
				log.Printf("controller unregistered (%v connected)", len(h.controllers))

			case s := <-h.RegisterScreen:
				h.screens[s] = true
				log.Printf("screen registered (%v connected)", len(h.screens))

			case s := <-h.UnregisterScreen:
				delete(h.screens, s)
				log.Printf("screen unregistered (%v connected)", len(h.screens))

			case c := <-h.Controls:
				if c.Type == "gameDone" {
					smtp, err := h.db.Prepare("INSERT INTO scores(playerID, playerName, finalScore, color) VALUES(?,?,?,?)")
					if err != nil {
						log.Fatalln(errors.Wrap(err, "error preparing create score DB statement"))
					}
					_, err = smtp.Exec(c.Player.ID, c.Player.Name, c.FinalScore, c.Color)
					if err != nil {
						log.Fatalln(errors.Wrap(err, "error executing create score DB statement"))
					}
				}

				for s := range h.screens {
					err := s.WriteJSON(c)
					if err != nil {
						err = s.Close()
						if err != nil {
							log.Fatalln(errors.Wrap(err, "error closing screen connection"))
						}
						h.UnregisterScreen <- s
					}
				}

			case e := <-h.Events:
				for c := range h.controllers {
					err := c.WriteJSON(e)
					if err != nil {
						err = c.Close()
						if err != nil {
							log.Fatalln(errors.Wrap(err, "error closing controller connection"))
						}
						h.UnregisterController <- c
					}
				}
			}
		}
	}()
}
