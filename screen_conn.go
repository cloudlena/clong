package clong

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

// ScreenConnHandler handles a WebSocket connection coming from a screen.
func ScreenConnHandler(h *Hub, up websocket.Upgrader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, errUpgradingConnection))
			return
		}
		defer ws.Close()

		h.RegisterScreen <- ws

		for {
			var e Event

			err = ws.ReadJSON(&e)
			if err != nil {
				handleHTTPError(w, errors.Wrap(err, errReadingJSON))
				h.UnregisterScreen <- ws
				break
			}

			h.Events <- e
		}
	})
}
