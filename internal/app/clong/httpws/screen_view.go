package httpws

import (
	"net/http"
	"path"
)

// HandleScreenView returns the screen HTML page.
func HandleScreenView() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("web", "static", "screen.html"))
	}
}
