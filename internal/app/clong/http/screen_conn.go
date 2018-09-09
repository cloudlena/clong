package clong

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(hub *Hub, up websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error upgrading connection"))
			return
		}
		defer func() {
			err := ws.Close()
			if err != nil {
				log.Fatal(errors.Wrap(err, "error closing websocket"))
			}
		}()

		hub.registerScreen <- ws

		for {
			var e event

			err = ws.ReadJSON(&e)
			if err != nil {
				handleHTTPError(w, errors.Wrap(err, "error reading JSON"))
				hub.unregisterScreen <- ws
				break
			}

			hub.events <- e
		}
	}
}
