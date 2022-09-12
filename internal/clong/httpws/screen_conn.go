package httpws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/gorilla/websocket"
)

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(svc clong.Service, up websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		conn, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error upgrading connection: %w", err))
			return
		}
		defer func() {
			err = conn.Close()
			if err != nil {
				log.Fatal(fmt.Errorf("error closing websocket: %w", err))
			}
		}()
		svc.RegisterScreen(conn)

		for {
			var evt clong.Event
			err = conn.ReadJSON(&evt)
			if err != nil {
				handleHTTPError(w, fmt.Errorf("error reading JSON: %w", err))
				svc.UnregisterScreen(conn)
				break
			}

			svc.PublishEvent(ctx, evt)
		}
	}
}
