package clong

import (
	"errors"
	"log"
	"net/http"
)

// Error codes returned by the application.
var (
	ErrUpgradingConnection = errors.New("error upgrading connection")
	ErrReadingMessage      = errors.New("error reading message")
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
