package httpws

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// maxMessageSize is the maximum size in bytes of a message read from a client.
	maxMessageSize = 4096
	// writeTimeout is how long writing a message to a client may take.
	writeTimeout = time.Second
)

// Clients are pinged regularly, so dead connections are noticed
// and proxies don't close idle ones (e.g. ingress-nginx after 60s).
// Variables, so tests can shorten them.
var (
	// pingPeriod is how often clients are pinged.
	pingPeriod = 15 * time.Second
	// pongWait is how long a client may take to answer a ping before its reads fail.
	pongWait = 30 * time.Second
)

var upgrader websocket.Upgrader

// clientConn is a WebSocket connection to a client.
// Its writes time out, so a client that stopped receiving can't block the service,
// which writes while holding its lock.
type clientConn struct {
	*websocket.Conn
}

// WriteJSON writes a message as JSON, giving up after writeTimeout.
func (c clientConn) WriteJSON(v any) error {
	if err := c.SetWriteDeadline(time.Now().Add(writeTimeout)); err != nil {
		return err
	}
	return c.Conn.WriteJSON(v)
}

// upgrade upgrades an HTTP connection to a WebSocket connection,
// which is pinged until the handler returns.
// It responds with an HTTP error itself if upgrading fails.
func upgrade(w http.ResponseWriter, r *http.Request) (clientConn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return clientConn{}, err
	}
	conn.SetReadLimit(maxMessageSize)
	if err := conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		closeConn(clientConn{conn})
		return clientConn{}, err
	}
	conn.SetPongHandler(func(string) error {
		return conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	// The request context is canceled when the handler returns
	go ping(r.Context(), conn)
	return clientConn{conn}, nil
}

// ping pings a connection every pingPeriod until the context is canceled or pinging fails.
func ping(ctx context.Context, conn *websocket.Conn) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			// WriteControl may be called concurrently with other writes
			if err := conn.WriteControl(websocket.PingMessage, nil, time.Now().Add(writeTimeout)); err != nil {
				return
			}
		}
	}
}

// closeConn closes a WebSocket connection and logs any error.
// Connections the service already closed are ignored.
func closeConn(conn clientConn) {
	if err := conn.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		log.Printf("error closing websocket connection: %v", err)
	}
}
