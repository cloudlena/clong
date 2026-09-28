package clong

import (
	"context"
	"errors"
	"log"
	"regexp"
	"sync"
	"time"
)

const (
	// gameDuration is how long a game lasts on the controller.
	gameDuration = 60 * time.Second
	// gameTolerance accounts for network latency and timer drift on the controller.
	gameTolerance = 3 * time.Second
	// staleGameAge is the age after which unfinished games are discarded.
	staleGameAge = 10 * time.Minute
)

// ErrScreenConnected is returned when a screen registers while another one is connected.
var ErrScreenConnected = errors.New("another screen is already connected")

var colorPattern = regexp.MustCompile(`^#[0-9a-fA-F]{6}$`)

// ClientConnection is a connection with a client.
type ClientConnection interface {
	WriteJSON(v any) error
	Close() error
}

// game is a game in progress, whose points are tallied by the server.
type game struct {
	started time.Time
	points  int64
}

// Service is a messaging hub between controllers, screens and scoreboards.
type Service struct {
	mu          sync.Mutex
	controllers map[ClientConnection]bool
	screens     map[ClientConnection]bool // holds at most one screen
	scoreboards map[ClientConnection]bool
	games       map[string]*game // by player ID
	scores      ScoreStore
}

// NewService creates a new service.
func NewService(scores ScoreStore) *Service {
	return &Service{
		controllers: make(map[ClientConnection]bool),
		screens:     make(map[ClientConnection]bool),
		scoreboards: make(map[ClientConnection]bool),
		games:       make(map[string]*game),
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

// RegisterScreen registers a new screen. Only one screen can be connected at a time.
func (s *Service) RegisterScreen(c ClientConnection) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.screens) > 0 {
		return ErrScreenConnected
	}
	s.screens[c] = true
	return nil
}

// UnregisterScreen removes a screen.
func (s *Service) UnregisterScreen(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.screens, c)
}

// RegisterScoreboard registers a new scoreboard.
func (s *Service) RegisterScoreboard(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.scoreboards[c] = true
}

// UnregisterScoreboard removes a scoreboard.
func (s *Service) UnregisterScoreboard(c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.scoreboards, c)
}

// PublishEvent sends an event to all controllers and tallies the points of running games.
func (s *Service) PublishEvent(event Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.Type == "BALL_DONE" {
		g, ok := s.games[event.Player.ID]
		if ok && time.Since(g.started) <= gameDuration+gameTolerance {
			g.points += event.Points
		}
	}
	broadcast(s.controllers, event)
}

// PublishControl handles a control from a controller.
// Game controls are handled by the service, all others are sent to the screen.
func (s *Service) PublishControl(ctx context.Context, ctrl Control) {
	switch ctrl.Type {
	case "GAME_STARTED":
		s.startGame(ctrl.Player)
	case "GAME_FINISHED":
		s.finishGame(ctx, ctrl)
	default:
		s.mu.Lock()
		defer s.mu.Unlock()
		broadcast(s.screens, ctrl)
	}
}

// startGame starts tallying points for a player.
func (s *Service) startGame(player User) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Discard games that were never finished, e.g. because the controller was closed
	for id, g := range s.games {
		if time.Since(g.started) > staleGameAge {
			delete(s.games, id)
		}
	}

	s.games[player.ID] = &game{started: time.Now()}
}

// finishGame saves the tallied score of a player's game and sends it to all scoreboards.
func (s *Service) finishGame(ctx context.Context, ctrl Control) {
	s.mu.Lock()
	g, ok := s.games[ctrl.Player.ID]
	delete(s.games, ctrl.Player.ID)
	s.mu.Unlock()

	if !ok {
		log.Printf("ignoring finished game of %s: game was never started\n", ctrl.Player.ID)
		return
	}
	if time.Since(g.started) < gameDuration-gameTolerance {
		log.Printf("ignoring finished game of %s: game finished too early\n", ctrl.Player.ID)
		return
	}
	if !colorPattern.MatchString(ctrl.Color) {
		log.Printf("ignoring finished game of %s: invalid color %q\n", ctrl.Player.ID, ctrl.Color)
		return
	}

	scr := Score{
		Player:     ctrl.Player,
		FinalScore: g.points,
		Color:      ctrl.Color,
	}
	if err := s.scores.Add(ctx, &scr); err != nil {
		log.Printf("error adding score to store: %v\n", err)
		return
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	broadcast(s.scoreboards, scr)
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
