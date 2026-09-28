package httpws

import (
	"log"
	"net/http"
	"net/url"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleControllerConn handles a WebSocket connection from a controller.
func HandleControllerConn(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userID, err := r.Cookie("userid")
		if err != nil {
			http.Error(w, "user ID missing", http.StatusUnauthorized)
			return
		}
		userName, err := r.Cookie("username")
		if err != nil {
			http.Error(w, "username missing", http.StatusUnauthorized)
			return
		}
		// Cookie values are URI-encoded by the controller
		name, err := url.PathUnescape(userName.Value)
		if err != nil {
			http.Error(w, "username invalid", http.StatusBadRequest)
			return
		}
		player := clong.User{ID: userID.Value, Name: name}
		if err := player.Validate(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// The upgrader responds with an HTTP error itself if upgrading fails
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			return
		}
		defer closeConn(conn)

		svc.RegisterController(conn)
		defer svc.UnregisterController(conn)

		for {
			var ctrl clong.Control
			if err := conn.ReadJSON(&ctrl); err != nil {
				log.Printf("error reading from controller: %v\n", err)
				return
			}
			ctrl.Player = player
			svc.PublishControl(r.Context(), ctrl)
		}
	}
}
