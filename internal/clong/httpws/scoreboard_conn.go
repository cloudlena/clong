package httpws

import (
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleScoreboardConn handles a WebSocket connection coming from a scoreboard.
func HandleScoreboardConn(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// The upgrader responds with an HTTP error itself if upgrading fails
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer closeConn(conn)

		svc.RegisterScoreboard(conn)
		defer svc.UnregisterScoreboard(conn)

		// Scoreboards only receive, but reading is needed to notice when they disconnect
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("error reading from scoreboard: %v\n", err)
				return
			}
		}
	}
}
