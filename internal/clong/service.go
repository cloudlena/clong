package clong

import (
	"context"
	"log"
	"sync"
)

// ClientConnection is a connection with a client.
type ClientConnection interface {
	WriteJSON(v any) error
	Close() error
}

// Service is a messaging hub between controllers and screens.
type Service struct {
	mu          sync.Mutex
	controllers map[ClientConnection]bool
	screens     map[ClientConnection]bool
	scores      ScoreStore
}

// NewService creates a new service.
func NewService(scores ScoreStore) *Service {
	return &Service{
		controllers: make(map[ClientConnection]bool),
		screens:     make(map[ClientConnection]bool),
		scores:      scores,
	}
}

// RegisterController registers a new controller.
func (s *Service) RegisterController(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.controllers[c] = true
}

// UnregisterController removes a controller.
func (s *Service) UnregisterController(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.controllers, c)
}

// RegisterScreen registers a new screen.
func (s *Service) RegisterScreen(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.screens[c] = true
}

// UnregisterScreen removes a screen.
func (s *Service) UnregisterScreen(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.screens, c)
}

// PublishEvent sends an event to all controllers.
func (s *Service) PublishEvent(event Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	broadcast(s.controllers, event)
}

// PublishControl sends a control to all screens and saves the score if the game is finished.
func (s *Service) PublishControl(ctx context.Context, ctrl Control) {
	if ctrl.Type == "GAME_FINISHED" {
		scr := Score{
			Player:     ctrl.Player,
			FinalScore: ctrl.FinalScore,
			Color:      ctrl.Color,
		}
		if err := s.scores.Add(ctx, &scr); err != nil {
			log.Printf("error adding score to store: %v\n", err)
		}
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	broadcast(s.screens, ctrl)
}

// broadcast sends a message to all connections and removes the ones that fail.
func broadcast(conns map[ClientConnection]bool, msg any) {
	for c := range conns {
		if err := c.WriteJSON(msg); err != nil {
			if err := c.Close(); err != nil {
				log.Printf("error closing connection: %v\n", err)
			}
			delete(conns, c)
		}
	}
}
