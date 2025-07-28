package clong_test

import (
	"context"
	"errors"
	"testing"

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
	added   []*clong.Score
	listAll []*clong.Score
	listErr error
	addErr  error
	rmErr   error
}

func (m *mockScoreStore) ListAll(_ context.Context) ([]*clong.Score, error) {
	return m.listAll, m.listErr
}

func (m *mockScoreStore) Add(_ context.Context, s *clong.Score) error {
	if m.addErr != nil {
		return m.addErr
	}
	m.added = append(m.added, s)
	return nil
}

func (m *mockScoreStore) RemoveAll(_ context.Context) error {
	return m.rmErr
}

func TestRegisterAndUnregisterController(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	conn := &mockConn{}

	svc.RegisterController(conn)
	svc.PublishEvent(context.Background(), clong.Event{Type: "ping"})
	if len(conn.written) != 1 {
		t.Fatalf("expected 1 message after register, got %d", len(conn.written))
	}

	svc.UnregisterController(conn)
	svc.PublishEvent(context.Background(), clong.Event{Type: "ping"})
	if len(conn.written) != 1 {
		t.Fatalf("expected no new message after unregister, got %d", len(conn.written))
	}
}

func TestRegisterAndUnregisterScreen(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	conn := &mockConn{}

	svc.RegisterScreen(conn)
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
	svc.PublishEvent(context.Background(), evt)

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
	svc.PublishEvent(context.Background(), clong.Event{Type: "test"})

	if !broken.closed {
		t.Error("expected broken controller to be closed")
	}

	// Broken conn is removed; a second publish should only reach healthy.
	svc.PublishEvent(context.Background(), clong.Event{Type: "test2"})
	if len(healthy.written) != 2 {
		t.Errorf("expected 2 messages on healthy conn, got %d", len(healthy.written))
	}
}

func TestPublishControlBroadcastsToAllScreens(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	s1, s2 := &mockConn{}, &mockConn{}
	svc.RegisterScreen(s1)
	svc.RegisterScreen(s2)

	svc.PublishControl(context.Background(), clong.Control{Type: "MOVE"})

	if len(s1.written) != 1 {
		t.Errorf("s1: expected 1 message, got %d", len(s1.written))
	}
	if len(s2.written) != 1 {
		t.Errorf("s2: expected 1 message, got %d", len(s2.written))
	}
}

func TestPublishControlRemovesBrokenScreen(t *testing.T) {
	svc := clong.NewService(&mockScoreStore{})
	broken := &mockConn{writeErr: errors.New("write failed")}
	svc.RegisterScreen(broken)

	svc.PublishControl(context.Background(), clong.Control{Type: "MOVE"})

	if !broken.closed {
		t.Error("expected broken screen to be closed")
	}
}

func TestPublishControlGameFinishedSavesScore(t *testing.T) {
	store := &mockScoreStore{}
	svc := clong.NewService(store)

	ctrl := clong.Control{
		Type:       "GAME_FINISHED",
		Player:     clong.User{ID: "u1", Name: "Alice"},
		FinalScore: 42,
		Color:      "red",
	}
	svc.PublishControl(context.Background(), ctrl)

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
	if got.Color != "red" {
		t.Errorf("Color: expected red, got %s", got.Color)
	}
}

func TestPublishControlNonGameFinishedDoesNotSaveScore(t *testing.T) {
	store := &mockScoreStore{}
	svc := clong.NewService(store)

	svc.PublishControl(context.Background(), clong.Control{Type: "MOVE", FinalScore: 99})

	if len(store.added) != 0 {
		t.Errorf("expected no score saved for MOVE, got %d", len(store.added))
	}
}
