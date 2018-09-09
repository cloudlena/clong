package httpws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// HandleControllerConn handles a WebSocket connection from a controller.
func HandleControllerConn(svc clong.Service, up websocket.Upgrader) http.HandlerFunc {
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
		svc.RegisterController(ws)

		for {
			userID, ok := cookieVal(r.Cookies(), "userid")
			if !ok {
				err := NewUnauthorizedError("user ID missing")
				handleHTTPError(w, err)
				svc.UnregisterController(ws)
				break
			}
			userName, ok := cookieVal(r.Cookies(), "username")
			if !ok {
				err := NewUnauthorizedError("user name missing")
				handleHTTPError(w, err)
				svc.UnregisterController(ws)
				break
			}

			var ctrl clong.Control
			err = ws.ReadJSON(&ctrl)
			if err != nil {
				handleHTTPError(w, errors.Wrap(err, "error reading JSON"))
				svc.UnregisterController(ws)
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
