package httpws

import (
	"errors"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/gorilla/websocket"
)

func TestHandleScreenConnRejectsSecondScreen(t *testing.T) {
	srv := httptest.NewServer(HandleScreenConn(clong.NewService(&mockScoreStore{})))
	defer srv.Close()

	dial(t, wsURL(srv.URL))
	second := dial(t, wsURL(srv.URL))

	_, _, err := second.ReadMessage()
	var closeErr *websocket.CloseError
	if !errors.As(err, &closeErr) || closeErr.Code != websocket.CloseTryAgainLater {
		t.Fatalf("expected close with code %d, got %v", websocket.CloseTryAgainLater, err)
	}
}

func TestHandleScreenConnRejectsOversizedMessages(t *testing.T) {
	srv := httptest.NewServer(HandleScreenConn(clong.NewService(&mockScoreStore{})))
	defer srv.Close()
	conn := dial(t, wsURL(srv.URL))

	// A valid event, so the connection can only be closed because of its size
	evt := `{"type":"` + strings.Repeat("a", 8192) + `"}`
	if err := conn.WriteMessage(websocket.TextMessage, []byte(evt)); err != nil {
		t.Fatalf("error writing message: %v", err)
	}

	if err := conn.SetReadDeadline(time.Now().Add(time.Second)); err != nil {
		t.Fatalf("error setting read deadline: %v", err)
	}
	_, _, err := conn.ReadMessage()
	var closeErr *websocket.CloseError
	if !errors.As(err, &closeErr) || closeErr.Code != websocket.CloseMessageTooBig {
		t.Fatalf("expected close with code %d, got %v", websocket.CloseMessageTooBig, err)
	}
}
