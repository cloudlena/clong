package clong

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ScreenConnHandler handles a WebSocket connections coming from a screen.
func ScreenConnHandler(h *Hub, up websocket.Upgrader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err, ErrUpgradingConnection)
			return
		}
		defer func() {
			err = ws.Close()
			if err != nil {
				log.Fatalln(err)
			}
		}()

		h.RegisterScreen <- ws

		for {
			var e Event

			err := ws.ReadJSON(&e)
			if err != nil {
				log.Println(ErrReadingMessage, err.Error())
				h.UnregisterScreen <- ws
				break
			}

			h.Events <- e
		}
	})
}
