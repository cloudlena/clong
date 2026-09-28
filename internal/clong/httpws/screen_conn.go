package httpws

import (
	"log"
	"net/http"
	"time"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/gorilla/websocket"
)

const closeTimeout = time.Second

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The upgrader responds with an HTTP error itself if upgrading fails
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer closeConn(conn)

		// Rejected screens are told why, since browsers can't read the status of a failed upgrade
		if err := svc.RegisterScreen(conn); err != nil {
			msg := websocket.FormatCloseMessage(websocket.CloseTryAgainLater, err.Error())
			if err := conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(closeTimeout)); err != nil {
				log.Printf("error rejecting screen: %v\n", err)
			}
			return
		}
		defer svc.UnregisterScreen(conn)

		for {
			var evt clong.Event
			if err := conn.ReadJSON(&evt); err != nil {
				log.Printf("error reading from screen: %v\n", err)
				return
			}
			svc.PublishEvent(evt)
		}
	}
}
