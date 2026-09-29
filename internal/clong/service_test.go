package clong_test

import (
	"context"
	"errors"
	"testing"
	"testing/synctest"
	"time"

	"github.com/cloudlena/clong/internal/clong"
)

// mockConn is a fake ClientConnection for testing.
type mockConn struct {
	written  []any
	writeErr error
	closed   bool
}

func (m *mockConn) WriteJSON(v any) error {
	if m.writeErr != nil {
		return m.writeErr
	}
	m.written = append(m.written, v)
	return nil
}

func (m *mockConn) Close() error {
	m.closed = true
	return nil
}

// mockScoreStore is a fake ScoreStore for testing.
type mockScoreStore struct {
	added []*clong.Score
}

func (m *mockScoreStore) ListAll(_ context.Context) ([]*clong.Score, error) {
	return m.added, nil
}

func (m *mockScoreStore) Add(_ context.Context, s *clong.Score) error {
	m.added = append(m.added, s)
	return nil
}

func (m *mockScoreStore) RemoveAll(_ context.Context) error {
	m.added = nil
	return nil
}

func TestRegisterAndUnregisterController(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	conn := &mockConn{}

	svc.RegisterController(conn)
	svc.PublishEvent(clong.Event{Type: "ping"})
	if len(conn.written) != 1 {
		t.Fatalf("expected 1 message after register, got %d", len(conn.written))
	}

	svc.UnregisterController(conn)
	svc.PublishEvent(clong.Event{Type: "ping"})
	if len(conn.written) != 1 {
		t.Fatalf("expected no new message after unregister, got %d", len(conn.written))
	}
}

func TestRegisterAndUnregisterScreen(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	conn := &mockConn{}

	if err := svc.RegisterScreen(conn); err != nil {
		t.Fatalf("unexpected error registering screen: %v", err)
	}
	svc.PublishControl(context.Background(), clong.Control{Type: "ping"})
	if len(conn.written) != 1 {
		t.Fatalf("expected 1 message after register, got %d", len(conn.written))
	}

	svc.UnregisterScreen(conn)
	svc.PublishControl(context.Background(), clong.Control{Type: "ping"})
	if len(conn.written) != 1 {
		t.Fatalf("expected no new message after unregister, got %d", len(conn.written))
	}
}

func TestPublishEventBroadcastsToAllControllers(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	c1, c2 := &mockConn{}, &mockConn{}
	svc.RegisterController(c1)
	svc.RegisterController(c2)

	evt := clong.Event{Type: "test", Points: 5}
	svc.PublishEvent(evt)

	if len(c1.written) != 1 {
		t.Errorf("c1: expected 1 message, got %d", len(c1.written))
	}
	if len(c2.written) != 1 {
		t.Errorf("c2: expected 1 message, got %d", len(c2.written))
	}
}

func TestPublishEventRemovesBrokenController(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	broken := &mockConn{writeErr: errors.New("write failed")}
	healthy := &mockConn{}

	svc.RegisterController(broken)
	svc.RegisterController(healthy)
	svc.PublishEvent(clong.Event{Type: "test"})

	if !broken.closed {
		t.Error("expected broken controller to be closed")
	}

	// Broken conn is removed; a second publish should only reach healthy.
	svc.PublishEvent(clong.Event{Type: "test2"})
	if len(healthy.written) != 2 {
		t.Errorf("expected 2 messages on healthy conn, got %d", len(healthy.written))
	}
}

func TestRegisterScreenRejectsSecondScreen(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	first, second := &mockConn{}, &mockConn{}

	if err := svc.RegisterScreen(first); err != nil {
		t.Fatalf("unexpected error registering first screen: %v", err)
	}
	if err := svc.RegisterScreen(second); !errors.Is(err, clong.ErrScreenConnected) {
		t.Fatalf("expected ErrScreenConnected, got %v", err)
	}

	// A new screen can connect once the first one is gone.
	svc.UnregisterScreen(first)
	if err := svc.RegisterScreen(second); err != nil {
		t.Fatalf("unexpected error registering screen after unregister: %v", err)
	}
}

func TestPublishControlRemovesBrokenScreen(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	broken := &mockConn{writeErr: errors.New("write failed")}
	if err := svc.RegisterScreen(broken); err != nil {
		t.Fatalf("unexpected error registering screen: %v", err)
	}

	svc.PublishControl(context.Background(), clong.Control{Type: "BALL_INIT"})

	if !broken.closed {
		t.Error("expected broken screen to be closed")
	}
	if err := svc.RegisterScreen(&mockConn{}); err != nil {
		t.Fatalf("expected broken screen to free its slot, got %v", err)
	}
}

// playGame starts a game for a player, scores the given points and finishes it after gameLength.
// It must be called inside a synctest bubble so time passes instantly.
func playGame(svc *clong.Service, player clong.User, color string, gameLength time.Duration, points ...int64) {
	svc.PublishControl(context.Background(), clong.Control{Type: "GAME_STARTED", Player: player})
	for _, p := range points {
		svc.PublishEvent(clong.Event{Type: "BALL_DONE", Player: player, Points: p})
	}
	time.Sleep(gameLength)
	svc.PublishControl(context.Background(), clong.Control{Type: "GAME_FINISHED", Player: player, Color: color})
}

func TestGameFinishedSavesTalliedScore(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		store := &mockScoreStore{}
		svc := clong.NewService(store)
		board := &mockConn{}
		svc.RegisterScoreboard(board)

		alice := clong.User{ID: "u1", Name: "Alice"}
		playGame(svc, alice, "#ff0000", 60*time.Second, 10, 0, 32)

		if len(store.added) != 1 {
			t.Fatalf("expected 1 saved score, got %d", len(store.added))
		}
		got := store.added[0]
		if got.FinalScore != 42 {
			t.Errorf("FinalScore: expected 42, got %d", got.FinalScore)
		}
		if got.Player.Name != "Alice" {
			t.Errorf("Player.Name: expected Alice, got %s", got.Player.Name)
		}
		if got.Color != "#ff0000" {
			t.Errorf("Color: expected #ff0000, got %s", got.Color)
		}
		if len(board.written) != 1 {
			t.Errorf("expected score to be sent to scoreboard, got %d messages", len(board.written))
		}
	})
}

func TestGameFinishedIgnoresPointsAfterGameEnded(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		store := &mockScoreStore{}
		svc := clong.NewService(store)
		alice := clong.User{ID: "u1", Name: "Alice"}

		svc.PublishControl(context.Background(), clong.Control{Type: "GAME_STARTED", Player: alice})
		svc.PublishEvent(clong.Event{Type: "BALL_DONE", Player: alice, Points: 10})
		time.Sleep(2 * time.Minute)
		svc.PublishEvent(clong.Event{Type: "BALL_DONE", Player: alice, Points: 50})
		svc.PublishControl(context.Background(), clong.Control{Type: "GAME_FINISHED", Player: alice, Color: "#ff0000"})

		if len(store.added) != 1 {
			t.Fatalf("expected 1 saved score, got %d", len(store.added))
		}
		if got := store.added[0].FinalScore; got != 10 {
			t.Errorf("FinalScore: expected 10, got %d", got)
		}
	})
}

func TestGameFinishedRejectsInvalidGames(t *testing.T) {
	alice := clong.User{ID: "u1", Name: "Alice"}
	tests := map[string]func(svc *clong.Service){
		"never started": func(svc *clong.Service) {
			svc.PublishControl(context.Background(), clong.Control{Type: "GAME_FINISHED", Player: alice, Color: "#ff0000"})
		},
		"finished too early": func(svc *clong.Service) {
			playGame(svc, alice, "#ff0000", 10*time.Second, 10)
		},
		"invalid color": func(svc *clong.Service) {
			playGame(svc, alice, "not-a-color", 60*time.Second, 10)
		},
	}

	for name, play := range tests {
		t.Run(name, func(t *testing.T) {
			synctest.Test(t, func(t *testing.T) {
				store := &mockScoreStore{}
				board := &mockConn{}
				svc := clong.NewService(store)
				svc.RegisterScoreboard(board)

				play(svc)

				if len(store.added) != 0 {
					t.Errorf("expected no saved score, got %d", len(store.added))
				}
				if len(board.written) != 0 {
					t.Errorf("expected nothing sent to scoreboard, got %d messages", len(board.written))
				}
			})
		})
	}
}

func TestGameControlsAreNotSentToScreen(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		svc := clong.NewService(&mockScoreStore{})
		screen := &mockConn{}
		if err := svc.RegisterScreen(screen); err != nil {
			t.Fatalf("unexpected error registering screen: %v", err)
		}

		playGame(svc, clong.User{ID: "u1", Name: "Alice"}, "#ff0000", 60*time.Second)

		if len(screen.written) != 0 {
			t.Errorf("expected no messages on screen, got %d", len(screen.written))
		}
	})
}

func TestPublishControlNonGameFinishedDoesNotSaveScore(t *testing.T) {
	store := &mockScoreStore{}
	svc := clong.NewService(store)

	svc.PublishControl(context.Background(), clong.Control{Type: "MOVE"})

	if len(store.added) != 0 {
		t.Errorf("expected no score saved for MOVE, got %d", len(store.added))
	}
}
