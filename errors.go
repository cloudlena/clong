package clong

import (
	"log"
	"net/http"
)

// Error codes commonly used throughout the application.
const (
	errUpgradingConnection = "error upgrading connection"
	errClosingConnection   = "error closing WebSocket connection"
	errReadingMessage      = "error reading message"
)

// handleHTTPError handles HTTP errors.
func handleHTTPError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError

	http.Error(w, http.StatusText(code), code)
	log.Println(err)
}
