package clong

import (
	"context"
	"log"
)

// ClientConnection is a connection with a client.
type ClientConnection interface {
	WriteJSON(v any) error
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

// BaseService is a messaging hub.
type BaseService struct {
	controllers map[ClientConnection]bool
	screens     map[ClientConnection]bool
	scores      ScoreStore
}

// NewService creates a new service.
func NewService(scores ScoreStore) *BaseService {
	svc := BaseService{
		controllers: make(map[ClientConnection]bool),
		screens:     make(map[ClientConnection]bool),
		scores:      scores,
	}
	return &svc
}

// RegisterController registers a new controller.
func (s *BaseService) RegisterController(c ClientConnection) {
	s.controllers[c] = true
}

// UnregisterController removes a controller.
func (s *BaseService) UnregisterController(c ClientConnection) {
	delete(s.controllers, c)
}

// RegisterScreen registers a new screen.
func (s *BaseService) RegisterScreen(c ClientConnection) {
	s.screens[c] = true
}

// UnregisterScreen removes a screen.
func (s *BaseService) UnregisterScreen(c ClientConnection) {
	delete(s.screens, c)
}

// PublishEvent publishes an event to the messaging bus.
func (s *BaseService) PublishEvent(_ context.Context, event Event) {
	for c := range s.controllers {
		err := c.WriteJSON(event)
		if err != nil {
			err = c.Close()
			if err != nil {
				log.Printf("error closing controller connection: %v\n", err)
			}
			s.UnregisterController(c)
		}
	}
}

// PublishControl publishes a control to the messaging bus.
func (s *BaseService) PublishControl(ctx context.Context, ctrl Control) {
	if ctrl.Type == "GAME_FINISHED" {
		scr := Score{
			Player:     ctrl.Player,
			FinalScore: ctrl.FinalScore,
			Color:      ctrl.Color,
		}
		err := s.scores.Add(ctx, &scr)
		if err != nil {
			log.Printf("error adding score to store: %v\n", err)
		}
	}

	for scrn := range s.screens {
		err := scrn.WriteJSON(ctrl)
		if err != nil {
			err = scrn.Close()
			if err != nil {
				log.Printf("error closing screen connection: %v\n", err)
			}
			s.UnregisterScreen(scrn)
		}
	}
}
