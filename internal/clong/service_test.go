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
	added     []clong.Score
	lastLimit int
}

func (m *mockScoreStore) ListTop(_ context.Context, limit int) ([]clong.Score, error) {
	m.lastLimit = limit
	return m.added, nil
}

func (m *mockScoreStore) Add(_ context.Context, s clong.Score) error {
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
	alice := clong.User{ID: "u1", Name: "Alice"}

	svc.RegisterController(alice.ID, conn)
	svc.PublishEvent(clong.Event{Type: "ping", Player: alice})
	if len(conn.written) != 1 {
		t.Fatalf("expected 1 message after register, got %d", len(conn.written))
	}

	svc.UnregisterController(alice.ID, conn)
	svc.PublishEvent(clong.Event{Type: "ping", Player: alice})
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

func TestPublishEventOnlyReachesControllersOfPlayer(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	alice1, alice2, bob := &mockConn{}, &mockConn{}, &mockConn{}
	svc.RegisterController("alice", alice1)
	svc.RegisterController("alice", alice2)
	svc.RegisterController("bob", bob)

	svc.PublishEvent(clong.Event{Type: "test", Player: clong.User{ID: "alice"}, Points: 5})

	if len(alice1.written) != 1 {
		t.Errorf("alice1: expected 1 message, got %d", len(alice1.written))
	}
	if len(alice2.written) != 1 {
		t.Errorf("alice2: expected 1 message, got %d", len(alice2.written))
	}
	if len(bob.written) != 0 {
		t.Errorf("bob: expected no message, got %d", len(bob.written))
	}
}

func TestPublishEventRemovesBrokenController(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	broken := &mockConn{writeErr: errors.New("write failed")}
	healthy := &mockConn{}

	alice := clong.User{ID: "u1", Name: "Alice"}

	svc.RegisterController(alice.ID, broken)
	svc.RegisterController(alice.ID, healthy)
	svc.PublishEvent(clong.Event{Type: "test", Player: alice})

	if !broken.closed {
		t.Error("expected broken controller to be closed")
	}

	// Broken conn is removed; a second publish should only reach healthy.
	svc.PublishEvent(clong.Event{Type: "test2", Player: alice})
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

func TestRegisterScoreboardSendsBestScores(t *testing.T) {
	store := &mockScoreStore{added: []clong.Score{{PlayerName: "Alice", FinalScore: 42, Color: "#ff0000"}}}
	svc := clong.NewService(store)
	board := &mockConn{}

	if err := svc.RegisterScoreboard(context.Background(), board); err != nil {
		t.Fatalf("unexpected error registering scoreboard: %v", err)
	}

	if store.lastLimit != 10 {
		t.Errorf("expected top 10 scores to be listed, got top %d", store.lastLimit)
	}
	if len(board.written) != 1 {
		t.Fatalf("expected 1 message, got %d", len(board.written))
	}
	if sent := board.written[0].([]clong.Score); len(sent) != 1 || sent[0] != store.added[0] {
		t.Errorf("expected scoreboard to receive %+v, got %+v", store.added, sent)
	}
}

func TestResetScoresClearsScoreboards(t *testing.T) {
	store := &mockScoreStore{added: []clong.Score{{PlayerName: "Alice", FinalScore: 42, Color: "#ff0000"}}}
	svc := clong.NewService(store)
	board := registerScoreboard(t, svc)

	if err := svc.ResetScores(context.Background()); err != nil {
		t.Fatalf("unexpected error resetting scores: %v", err)
	}

	if len(store.added) != 0 {
		t.Errorf("expected store to be empty, got %d scores", len(store.added))
	}
	if len(board.written) != 1 {
		t.Fatalf("expected 1 message, got %d", len(board.written))
	}
	if sent := board.written[0].([]clong.Score); len(sent) != 0 {
		t.Errorf("expected scoreboard to receive no scores, got %+v", sent)
	}
}

// registerScoreboard registers a scoreboard and discards the scores sent on registration.
func registerScoreboard(t *testing.T, svc *clong.Service) *mockConn {
	t.Helper()
	board := &mockConn{}
	if err := svc.RegisterScoreboard(context.Background(), board); err != nil {
		t.Fatalf("unexpected error registering scoreboard: %v", err)
	}
	board.written = nil
	return board
}

// playGame starts a game for a player, scores the given points and finishes it after gameLength.
// It must be called inside a synctest bubble so time passes instantly.
func playGame(svc *clong.Service, player clong.User, color string, gameLength time.Duration, points ...int64) {
	svc.PublishControl(context.Background(), clong.Control{Type: clong.ControlGameStarted, Player: player})
	for _, p := range points {
		svc.PublishEvent(clong.Event{Type: clong.EventBallDone, Player: player, Points: p})
	}
	time.Sleep(gameLength)
	svc.PublishControl(context.Background(), clong.Control{Type: clong.ControlGameFinished, Player: player, Color: color})
}

func TestGameFinishedSavesTalliedScore(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		store := &mockScoreStore{}
		svc := clong.NewService(store)
		board := registerScoreboard(t, svc)

		alice := clong.User{ID: "u1", Name: "Alice"}
		playGame(svc, alice, "#ff0000", 60*time.Second, 10, 0, 32)

		if len(store.added) != 1 {
			t.Fatalf("expected 1 saved score, got %d", len(store.added))
		}
		got := store.added[0]
		if got.FinalScore != 42 {
			t.Errorf("FinalScore: expected 42, got %d", got.FinalScore)
		}
		if got.PlayerID != "u1" || got.PlayerName != "Alice" {
			t.Errorf("Player: expected u1/Alice, got %s/%s", got.PlayerID, got.PlayerName)
		}
		if got.Color != "#ff0000" {
			t.Errorf("Color: expected #ff0000, got %s", got.Color)
		}
		if len(board.written) != 1 {
			t.Fatalf("expected scores to be sent to scoreboard, got %d messages", len(board.written))
		}
		if sent := board.written[0].([]clong.Score); len(sent) != 1 || sent[0] != got {
			t.Errorf("expected scoreboard to receive %+v, got %+v", got, sent)
		}
	})
}

func TestGameFinishedIgnoresPointsAfterGameEnded(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		store := &mockScoreStore{}
		svc := clong.NewService(store)
		alice := clong.User{ID: "u1", Name: "Alice"}

		svc.PublishControl(context.Background(), clong.Control{Type: clong.ControlGameStarted, Player: alice})
		svc.PublishEvent(clong.Event{Type: clong.EventBallDone, Player: alice, Points: 10})
		time.Sleep(2 * time.Minute)
		svc.PublishEvent(clong.Event{Type: clong.EventBallDone, Player: alice, Points: 50})
		svc.PublishControl(context.Background(), clong.Control{Type: clong.ControlGameFinished, Player: alice, Color: "#ff0000"})

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
			svc.PublishControl(context.Background(), clong.Control{Type: clong.ControlGameFinished, Player: alice, Color: "#ff0000"})
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
				svc := clong.NewService(store)
				board := registerScoreboard(t, svc)

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
