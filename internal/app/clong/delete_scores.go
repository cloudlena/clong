package clong

import (
	"net/http"

	"github.com/mastertinner/clong/internal/app/clong/scores"
	"github.com/pkg/errors"
)

// HandleDeleteScores deletes all scores and resets the scoreboard.
func HandleDeleteScores(repo scores.Repository) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Delete scores from DB
		ctx := r.Context()
		err := repo.RemoveAll(ctx)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error removing scores"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
