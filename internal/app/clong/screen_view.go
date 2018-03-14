package clong

import (
	"net/http"
	"path"
)

// ScreenViewHandler returns the screen HTML page.
func ScreenViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, path.Join("web", "static", "screen.html"))
	})
}
