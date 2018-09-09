package httpws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(svc clong.Service, up websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error upgrading connection"))
			return
		}
		defer func() {
			err = ws.Close()
			if err != nil {
				log.Fatal(errors.Wrap(err, "error closing websocket"))
			}
		}()
		svc.RegisterScreen(ws)

		for {
			var evt clong.Event
			err = ws.ReadJSON(&evt)
			if err != nil {
				handleHTTPError(w, errors.Wrap(err, "error reading JSON"))
				svc.UnregisterScreen(ws)
				break
			}

			svc.PublishEvent(ctx, evt)
		}
	}
}
