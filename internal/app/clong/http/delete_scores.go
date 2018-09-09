package clong

import (
	"net/http"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// HandleDeleteScores deletes all scores and resets the scoreboard.
func HandleDeleteScores(store clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Delete scores from DB
		ctx := r.Context()
		err := store.RemoveAll(ctx)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error removing scores"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
