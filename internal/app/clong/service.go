package clong

import (
	"context"
	"fmt"
	"log"
)

// ClientConnection is a connection with a client.
type ClientConnection interface {
	WriteJSON(v interface{}) error
	Close() error
}

// Service is a messaging hub for events.
type Service interface {
	RegisterScreen(screen ClientConnection)
	UnregisterScreen(screen ClientConnection)
	RegisterController(controller ClientConnection)
	UnregisterController(controller ClientConnection)
	PublishEvent(ctx context.Context, event Event)
	PublishControl(ctx context.Context, ctrl Control)
}

// service is a messaging hub.
type service struct {
	controllers map[ClientConnection]bool
	screens     map[ClientConnection]bool
	scores      ScoreStore
}

// NewService creates a new service.
func NewService(scores ScoreStore) Service {
	svc := service{
		controllers: make(map[ClientConnection]bool),
		screens:     make(map[ClientConnection]bool),
		scores:      scores,
	}
	return &svc
}

// RegisterController registers a new controller.
func (s *service) RegisterController(c ClientConnection) {
	s.controllers[c] = true
}

// UnregisterController removes a controller.
func (s *service) UnregisterController(c ClientConnection) {
	delete(s.controllers, c)
}

// RegisterScreen registers a new screen.
func (s *service) RegisterScreen(c ClientConnection) {
	s.screens[c] = true
}

// UnregisterScreen removes a screen.
func (s *service) UnregisterScreen(c ClientConnection) {
	delete(s.screens, c)
}

// EventsChannel returns the hub's events channel.
func (s *service) PublishEvent(_ context.Context, event Event) {
	for c := range s.controllers {
		err := c.WriteJSON(event)
		if err != nil {
			err = c.Close()
			if err != nil {
				log.Fatal(fmt.Errorf("error closing controller connection: %w", err))
			}
			s.UnregisterController(c)
		}
	}
}

// PublishControl returns the hub's events channel.
func (s *service) PublishControl(ctx context.Context, ctrl Control) {
	switch ctrl.Type {
	case "GAME_FINISHED":
		scr := Score{
			Player:     ctrl.Player,
			FinalScore: ctrl.FinalScore,
			Color:      ctrl.Color,
		}
		err := s.scores.Add(ctx, &scr)
		if err != nil {
			log.Fatal(fmt.Errorf("error adding score to store: %w", err))
		}
	default:
		for scrn := range s.screens {
			err := scrn.WriteJSON(ctrl)
			if err != nil {
				err = scrn.Close()
				if err != nil {
					log.Fatal(fmt.Errorf("error closing screen connection: %w", err))
				}
				s.UnregisterScreen(scrn)
			}
		}
	}
}
