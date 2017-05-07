package clong

import (
	"log"
	"net/http"
)

// These are common errors used throughout the application.
const (
	ErrUpgradingConnection = "error upgrading connection"
	ErrReadingMessage      = "error reading message"
)

// handleHTTPError handles HTTP errors.
func handleHTTPError(w http.ResponseWriter, statusCode int, err error, msg string) {
	extMsg := http.StatusText(statusCode)
	if msg != "" {
		extMsg = extMsg + ": " + msg
	}
	http.Error(w, extMsg, statusCode)

	logMsg := extMsg
	if err != nil && err.Error() != msg {
		logMsg = logMsg + ": " + err.Error()
	}
	log.Println(logMsg)
}
