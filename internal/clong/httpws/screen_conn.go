package httpws

import (
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleScreenConn handles a WebSocket connection coming from a screen.
func HandleScreenConn(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The upgrader responds with an HTTP error itself if upgrading fails
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer closeConn(conn)

		svc.RegisterScreen(conn)
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
