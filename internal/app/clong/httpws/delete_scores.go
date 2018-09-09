package httpws

import (
	"net/http"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// HandleDeleteScores deletes all scores and resets the scoreboard.
func HandleDeleteScores(scores clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := scores.RemoveAll(ctx)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error removing all scores from store"))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}
