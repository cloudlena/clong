package clong

import (
	"net/http"

	"github.com/pkg/errors"
)

// HandleDeleteScores deletes all scores and resets the scoreboard.
func HandleDeleteScores(db ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Delete scores from DB
		ctx := r.Context()
		err := db.DeleteScores(ctx)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error removing scores"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
