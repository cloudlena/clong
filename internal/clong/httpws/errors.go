package httpws

import (
	"log"
	"net/http"
)

// handleServerError logs an error and responds with an internal server error.
func handleServerError(w http.ResponseWriter, err error) {
	log.Println(err)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}
