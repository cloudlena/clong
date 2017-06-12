package clong

import (
	"fmt"
	"net/http"
	"path"
)

// ScreenViewHandler returns the screen HTML page.
func ScreenViewHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("TLS", r.TLS)
		http.ServeFile(w, r, path.Join("public", "screen.html"))
	})
}
