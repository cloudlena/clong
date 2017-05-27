package clong

import (
	"errors"
	"log"
	"net/http"
)

// Error codes commonly used throughout the application.
const (
	errUpgradingConnection = "error upgrading connection"
	errClosingConnection   = "error closing WebSocket connection"
	errReadingMessage      = "error reading message"
)

// Errors commonly used throughout the application.
var (
	errUserIDMissing = errors.New("user ID missing")
)

// handleHTTPError handles HTTP errors.
func handleHTTPError(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError

	if err == errUserIDMissing {
		code = http.StatusUnauthorized
	}

	http.Error(w, http.StatusText(code), code)

	// Log if server error
	if code >= 500 {
		log.Println(err)
	}
}
