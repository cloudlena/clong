package clong

import (
	"net/http"
	"path"
)

// ScoreboardViewHandler returns the screen HTML page.
func ScoreboardViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("web", "static", "scoreboard.html"))
	})
}
