package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// ControllerConnHandler handles a WebSocket connection coming from a controller
func ControllerConnHandler(h *Hub, up websocket.Upgrader) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := up.Upgrade(w, r, nil)
		if err != nil {
			handleHTTPError(w, http.StatusInternalServerError, err, ErrUpgradingConnection)
			return
		}
		defer ws.Close()

		h.RegisterController <- ws

		for {
			var c Control
			id, ok := cookieVal(r.Cookies(), "userid")
			if !ok {
				log.Println("no userID foundresult")
				h.UnregisterController <- ws
				break
			}
			c.Player = id

			err := ws.ReadJSON(&c)
			if err != nil {
				log.Println(ErrReadingMessage, err.Error())
				h.UnregisterController <- ws
				break
			}

			h.Controls <- c
		}
	})
}

// cookieVal returns the value of a cookie
func cookieVal(cookies []*http.Cookie, name string) (string, bool) {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value, true
		}
	}

	return "", false
}
