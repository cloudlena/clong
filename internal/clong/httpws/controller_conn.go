package httpws

import (
	"fmt"
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/gorilla/websocket"
)

// HandleControllerConn handles a WebSocket connection from a controller.
func HandleControllerConn(svc clong.Service, up websocket.Upgrader) http.HandlerFunc {
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
		svc.RegisterController(conn)

		for {
			userID, ok := cookieVal(r.Cookies(), "userid")
			if !ok {
				err := NewUnauthorizedError("user ID missing")
				handleHTTPError(w, err)
				svc.UnregisterController(conn)
				break
			}
			userName, ok := cookieVal(r.Cookies(), "username")
			if !ok {
				err := NewUnauthorizedError("username missing")
				handleHTTPError(w, err)
				svc.UnregisterController(conn)
				break
			}

			var ctrl clong.Control
			err = conn.ReadJSON(&ctrl)
			if err != nil {
				handleHTTPError(w, fmt.Errorf("error reading JSON: %w", err))
				svc.UnregisterController(conn)
				break
			}
			ctrl.Player = clong.User{
				ID:   userID,
				Name: userName,
			}

			svc.PublishControl(ctx, ctrl)
		}
	}
}

// cookieVal returns the value of a cookie.
func cookieVal(cookies []*http.Cookie, name string) (string, bool) {
	for _, c := range cookies {
		if c.Name == name {
			return c.Value, true
		}
	}
	return "", false
}
