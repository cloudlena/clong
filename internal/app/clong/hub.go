package clong

import (
	"context"
	"log"

	"github.com/gorilla/websocket"
	"github.com/mastertinner/clong/internal/app/clong/scores"
	"github.com/pkg/errors"
)

// Hub is a WebSocket messaging hub.
type Hub struct {
	registerController   chan *websocket.Conn
	unregisterController chan *websocket.Conn
	registerScreen       chan *websocket.Conn
	unregisterScreen     chan *websocket.Conn
	controls             chan control
	events               chan event
	controllers          map[*websocket.Conn]bool
	screens              map[*websocket.Conn]bool
	repo                 scores.Repository
}

// NewHub creates a new messaging hub.
func NewHub(repo scores.Repository) *Hub {
	return &Hub{
		registerController:   make(chan *websocket.Conn),
		unregisterController: make(chan *websocket.Conn),
		registerScreen:       make(chan *websocket.Conn),
		unregisterScreen:     make(chan *websocket.Conn),
		controls:             make(chan control),
		events:               make(chan event),
		controllers:          make(map[*websocket.Conn]bool),
		screens:              make(map[*websocket.Conn]bool),
		repo:                 repo,
	}
}

// Run runs the messaging hub in a forever loop as a goruotine.
func (h *Hub) Run() { // nolint: gocyclo
	go func() {
		for {
			select {
			case c := <-h.registerController:
				h.controllers[c] = true
				log.Printf("controller registered (%v connected)", len(h.controllers))

			case c := <-h.unregisterController:
				delete(h.controllers, c)
				log.Printf("controller unregistered (%v connected)", len(h.controllers))

			case s := <-h.registerScreen:
				h.screens[s] = true
				log.Printf("screen registered (%v connected)", len(h.screens))

			case s := <-h.unregisterScreen:
				delete(h.screens, s)
				log.Printf("screen unregistered (%v connected)", len(h.screens))

			case c := <-h.controls:
				if c.Type == "GAME_FINISHED" {
					s := scores.Score{
						Player:     c.Player,
						FinalScore: c.FinalScore,
						Color:      c.Color,
					}
					ctx := context.Background()
					err := h.repo.Add(ctx, s)
					if err != nil {
						log.Fatal(errors.Wrap(err, "error creating score in DB"))
					}
				}

				for s := range h.screens {
					err := s.WriteJSON(c)
					if err != nil {
						err = s.Close()
						if err != nil {
							log.Fatal(errors.Wrap(err, "error closing screen connection"))
						}
						h.unregisterScreen <- s
					}
				}

			case e := <-h.events:
				for c := range h.controllers {
					err := c.WriteJSON(e)
					if err != nil {
						err = c.Close()
						if err != nil {
							log.Fatal(errors.Wrap(err, "error closing controller connection"))
						}
						h.unregisterController <- c
					}
				}
			}
		}
	}()
}
