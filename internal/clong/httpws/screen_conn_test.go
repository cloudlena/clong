package httpws_test

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
	"github.com/gorilla/websocket"
)

func TestHandleScreenConnRejectsSecondScreen(t *testing.T) {
	srv := httptest.NewServer(httpws.HandleScreenConn(clong.NewService(&mockScoreStore{})))
	defer srv.Close()
	url := "ws" + strings.TrimPrefix(srv.URL, "http")

	dialScreen(t, url)
	second := dialScreen(t, url)

	_, _, err := second.ReadMessage()
	var closeErr *websocket.CloseError
	if !errors.As(err, &closeErr) || closeErr.Code != websocket.CloseTryAgainLater {
		t.Fatalf("expected close with code %d, got %v", websocket.CloseTryAgainLater, err)
	}
}

// dialScreen connects to a screen WebSocket and closes the connection when the test ends.
func dialScreen(t *testing.T, url string) *websocket.Conn {
	t.Helper()
	conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatalf("error connecting screen: %v", err)
	}
	t.Cleanup(func() {
		if err := conn.Close(); err != nil {
			t.Errorf("error closing screen connection: %v", err)
		}
	})
	return conn
}
