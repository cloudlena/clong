package httpws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mastertinner/clong/internal/app/clong"
)

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(svc clong.Service, up websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, fmt.Errorf("error upgrading connection: %w", err))
			return
		}
		defer func() {
			err = ws.Close()
			if err != nil {
				log.Fatal(fmt.Errorf("error closing websocket: %w", err))
			}
		}()
		svc.RegisterScreen(ws)

		for {
			var evt clong.Event
			err = ws.ReadJSON(&evt)
			if err != nil {
				handleHTTPError(w, fmt.Errorf("error reading JSON: %w", err))
				svc.UnregisterScreen(ws)
				break
			}

			svc.PublishEvent(ctx, evt)
		}
	}
}
