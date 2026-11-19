package httpws

import (
	"context"
	"strings"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/gorilla/websocket"
)

type mockScoreStore struct {
	scores []clong.Score
	rmErr  error
}

func (m *mockScoreStore) ListTop(_ context.Context, _ int) ([]clong.Score, error) {
	return m.scores, nil
}

func (m *mockScoreStore) Add(_ context.Context, s clong.Score) error {
	m.scores = append(m.scores, s)
	return nil
}

func (m *mockScoreStore) RemoveAll(_ context.Context) error {
	return m.rmErr
}

// wsURL turns the URL of an HTTP test server into a WebSocket URL.
func wsURL(httpURL string) string {
	return "ws" + strings.TrimPrefix(httpURL, "http")
}

// dial connects to a WebSocket server and closes the connection when the test ends.
func dial(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("error connecting: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("error closing connection: %v", err)
		}
	})
	return conn
}
