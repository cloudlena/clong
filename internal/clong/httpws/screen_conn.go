package httpws

import (
	"log"
	"net/http"
	"time"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/gorilla/websocket"
)

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrade(w, r)
		if err != nil {
			return
		}
		defer closeConn(conn)

		// Rejected screens are told why, since browsers can't read the status of a failed upgrade
		if err := svc.RegisterScreen(conn); err != nil {
			msg := websocket.FormatCloseMessage(websocket.CloseTryAgainLater, err.Error())
			if err := conn.WriteControl(websocket.CloseMessage, msg, time.Now().Add(writeTimeout)); err != nil {
				log.Printf("error rejecting screen: %v", err)
			}
			return
		}
		defer svc.UnregisterScreen(conn)

		for {
			var evt clong.Event
			if err := conn.ReadJSON(&evt); err != nil {
				log.Printf("error reading from screen: %v", err)
				return
			}
			svc.PublishEvent(evt)
		}
	}
}
