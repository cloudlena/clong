package clong

import (
	"net/http"
	"path"
)

// HandleScoreboardView returns the screen HTML page.
func HandleScoreboardView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("web", "static", "scoreboard.html"))
	}
}
