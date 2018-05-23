package clong

import (
	"net/http"

	"github.com/pkg/errors"
)

// RemoveScoresHandler removes all scores and resets the scoreboard.
func RemoveScoresHandler(db ScoreStore) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Remove scores from DB
		err := db.RemoveScores()
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error removing scores"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	})
}
