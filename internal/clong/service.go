package clong

import (
	"context"
	"errors"
	"fmt"
	"log"
	"regexp"
	"sync"
	"time"
)

const (
	// gameDuration is how long a game lasts on the controller (see gameSeconds in controller.js).
	gameDuration = 60 * time.Second
	// gameTolerance accounts for network latency and timer drift on the controller.
	gameTolerance = 3 * time.Second
	// staleGameAge is the age after which unfinished games are discarded.
	staleGameAge = 10 * time.Minute
	// maxScores is how many scores scoreboards show.
	maxScores = 10
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
	controllers map[string]map[ClientConnection]bool // by player ID
	screens     map[ClientConnection]bool            // holds at most one screen
	games       map[string]*game                     // by player ID

	// scoreboardsMu also serializes sending scores, so scoreboards end up with the latest ones.
	scoreboardsMu sync.Mutex
	scoreboards   map[ClientConnection]bool
	scores        ScoreStore
}

// NewService creates a new service.
func NewService(scores ScoreStore) *Service {
	return &Service{
		controllers: make(map[string]map[ClientConnection]bool),
		screens:     make(map[ClientConnection]bool),
		scoreboards: make(map[ClientConnection]bool),
		games:       make(map[string]*game),
		scores:      scores,
	}
}

// RegisterController registers a new controller of a player.
func (s *Service) RegisterController(playerID string, c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.controllers[playerID] == nil {
		s.controllers[playerID] = make(map[ClientConnection]bool)
	}
	s.controllers[playerID][c] = true
}

// UnregisterController removes a controller of a player.
func (s *Service) UnregisterController(playerID string, c ClientConnection) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.controllers[playerID], c)
	if len(s.controllers[playerID]) == 0 {
		delete(s.controllers, playerID)
	}
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

// RegisterScoreboard sends the best scores to a new scoreboard and registers it to receive updates.
func (s *Service) RegisterScoreboard(ctx context.Context, c ClientConnection) error {
	s.scoreboardsMu.Lock()
	defer s.scoreboardsMu.Unlock()
	scrs, err := s.scores.ListTop(ctx, maxScores)
	if err != nil {
		return fmt.Errorf("error listing scores: %w", err)
	}
	if err := c.WriteJSON(scrs); err != nil {
		return fmt.Errorf("error sending scores: %w", err)
	}
	s.scoreboards[c] = true
	return nil
}

// UnregisterScoreboard removes a scoreboard.
func (s *Service) UnregisterScoreboard(c ClientConnection) {
	s.scoreboardsMu.Lock()
	defer s.scoreboardsMu.Unlock()
	delete(s.scoreboards, c)
}

// ResetScores removes all scores and clears all scoreboards.
func (s *Service) ResetScores(ctx context.Context) error {
	if err := s.scores.RemoveAll(ctx); err != nil {
		return fmt.Errorf("error removing all scores from store: %w", err)
	}
	s.publishScores(ctx)
	return nil
}

// publishScores sends the best scores to all scoreboards.
func (s *Service) publishScores(ctx context.Context) {
	s.scoreboardsMu.Lock()
	defer s.scoreboardsMu.Unlock()
	scrs, err := s.scores.ListTop(ctx, maxScores)
	if err != nil {
		log.Printf("error listing scores: %v", err)
		return
	}
	broadcast(s.scoreboards, scrs)
}

// PublishEvent sends an event to the controllers of its player and tallies the points of running games.
// Events aren't sent to other players, so they can't learn each other's IDs and interfere with each other's games.
func (s *Service) PublishEvent(event Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if event.Type == EventBallDone {
		g, ok := s.games[event.Player.ID]
		if ok && time.Since(g.started) <= gameDuration+gameTolerance {
			g.points += event.Points
		}
	}
	broadcast(s.controllers[event.Player.ID], event)
}

// PublishControl handles a control from a controller.
// Game controls are handled by the service, all others are sent to the screen.
func (s *Service) PublishControl(ctx context.Context, ctrl Control) {
	switch ctrl.Type {
	case ControlGameStarted:
		s.startGame(ctrl.Player)
	case ControlGameFinished:
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

// finishGame saves the tallied score of a player's game and updates all scoreboards.
func (s *Service) finishGame(ctx context.Context, ctrl Control) {
	s.mu.Lock()
	g, ok := s.games[ctrl.Player.ID]
	delete(s.games, ctrl.Player.ID)
	s.mu.Unlock()

	if !ok {
		log.Printf("ignoring finished game of %s: game was never started", ctrl.Player.ID)
		return
	}
	if time.Since(g.started) < gameDuration-gameTolerance {
		log.Printf("ignoring finished game of %s: game finished too early", ctrl.Player.ID)
		return
	}
	if !colorPattern.MatchString(ctrl.Color) {
		log.Printf("ignoring finished game of %s: invalid color %q", ctrl.Player.ID, ctrl.Color)
		return
	}

	scr := Score{
		PlayerID:   ctrl.Player.ID,
		PlayerName: ctrl.Player.Name,
		FinalScore: g.points,
		Color:      ctrl.Color,
	}
	if err := s.scores.Add(ctx, scr); err != nil {
		log.Printf("error adding score to store: %v", err)
		return
	}
	s.publishScores(ctx)
}

// broadcast sends a message to all connections and removes the ones that fail.
func broadcast(conns map[ClientConnection]bool, msg any) {
	for c := range conns {
		if err := c.WriteJSON(msg); err != nil {
			if err := c.Close(); err != nil {
				log.Printf("error closing connection: %v", err)
			}
			delete(conns, c)
		}
	}
}
