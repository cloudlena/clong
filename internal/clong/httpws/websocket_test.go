package httpws

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestClientConnWriteTimesOut(t *testing.T) {
	conns := make(chan clientConn, 1)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrade(w, r)
		if err != nil {
			t.Errorf("error upgrading connection: %v", err)
			return
		}
		conns <- conn
	}))
	defer srv.Close()

	// The client never reads, so the server's writes block once the network buffers are full
	dial(t, wsURL(srv.URL))
	conn := <-conns
	defer closeConn(conn)

	msg := strings.Repeat("a", 1<<16)
	for range 10000 {
		start := time.Now()
		if err := conn.WriteJSON(msg); err != nil {
			if d := time.Since(start); d > 2*writeTimeout {
				t.Fatalf("expected write to time out after %v, took %v", writeTimeout, d)
			}
			return
		}
	}
	t.Fatal("expected write to time out")
}

func TestUpgradeDropsClientsThatDontAnswerPings(t *testing.T) {
	shortenKeepalive(t)
	readErrs := serveReads(t)

	// The client never reads, so it never answers pings
	dial(t, readErrs.url)

	select {
	case err := <-readErrs.ch:
		t.Logf("read failed as expected: %v", err)
	case <-time.After(10 * pongWait):
		t.Fatal("expected read to fail when client doesn't answer pings")
	}
}

func TestUpgradeKeepsClientsThatAnswerPings(t *testing.T) {
	shortenKeepalive(t)
	readErrs := serveReads(t)

	// Reading makes the client answer pings
	client := dial(t, readErrs.url)
	go func() {
		for {
			if _, _, err := client.ReadMessage(); err != nil {
				return
			}
		}
	}()

	select {
	case err := <-readErrs.ch:
		t.Fatalf("expected connection to stay open, got %v", err)
	case <-time.After(5 * pongWait):
	}
}

// shortenKeepalive makes pings frequent for the duration of a test.
func shortenKeepalive(t *testing.T) {
	t.Helper()
	oldPingPeriod, oldPongWait := pingPeriod, pongWait
	pingPeriod, pongWait = 20*time.Millisecond, 100*time.Millisecond
	t.Cleanup(func() { pingPeriod, pongWait = oldPingPeriod, oldPongWait })
}

type readServer struct {
	url string
	ch  chan error
}

// serveReads starts a server that reads from WebSocket connections and reports when reading fails.
func serveReads(t *testing.T) readServer {
	t.Helper()
	ch := make(chan error, 1)
	var handlers sync.WaitGroup
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handlers.Add(1)
		defer handlers.Done()
		conn, err := upgrade(w, r)
		if err != nil {
			t.Errorf("error upgrading connection: %v", err)
			return
		}
		defer closeConn(conn)
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				ch <- err
				return
			}
		}
	}))
	// Closing the server doesn't wait for hijacked connections
	t.Cleanup(func() {
		srv.Close()
		handlers.Wait()
	})
	return readServer{url: wsURL(srv.URL), ch: ch}
}
