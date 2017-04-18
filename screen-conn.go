package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ScreenConnHandler handles a WebSocket connections coming from a screen
func ScreenConnHandler(h *Hub, up websocket.Upgrader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err, ErrUpgradingConnection)
			return
		}
		defer ws.Close()

		h.RegisterScreen <- ws

		for {
			var msg Message

			err := ws.ReadJSON(&msg)
			if err != nil {
				log.Println(ErrReadingMessage, err.Error())
				h.UnregisterScreen <- ws
				break
			}
		}
	})
}
