package clong

import (
	"log"
	"net/http"

	"github.com/pkg/errors"
)

// RemoveScoresHandler removes all scores and resets the scoreboard.
func RemoveScoresHandler(db DB) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove scores from DB
		err := removeScores(db)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error removing scores"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}

// removeScores removes all scores from the DB.
func removeScores(db DB) error {
	rows, err := db.Query("DELETE FROM scores")
	if err != nil {
		return errors.Wrap(err, "error removing scores from DB")
	}
	defer func() {
		err := rows.Close()
		if err != nil {
			log.Fatal(errors.Wrap(err, "error closing DB rows"))
		}
	}()

	return nil
}
