package httpws

import (
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleScoreboardConn handles a WebSocket connection coming from a scoreboard.
func HandleScoreboardConn(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrade(w, r)
		if err != nil {
			return
		}
		defer closeConn(conn)

		if err := svc.RegisterScoreboard(r.Context(), conn); err != nil {
			log.Printf("error registering scoreboard: %v", err)
			return
		}
		defer svc.UnregisterScoreboard(conn)

		// Scoreboards only receive, but reading is needed to notice when they disconnect
		for {
			if _, _, err := conn.ReadMessage(); err != nil {
				log.Printf("error reading from scoreboard: %v", err)
				return
			}
		}
	}
}
