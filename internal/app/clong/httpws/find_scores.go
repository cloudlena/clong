package httpws

import (
	"encoding/json"
	"net/http"

	"github.com/mastertinner/clong/internal/app/clong"
	"github.com/pkg/errors"
)

// HandleFindScores returns all scores as JSON.
func HandleFindScores(scores clong.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		scrs, err := scores.ListAll(ctx)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error finding scores"))
			return
		}

		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		err = json.NewEncoder(w).Encode(scrs)
		if err != nil {
			handleHTTPError(w, errors.Wrap(err, "error encoding JSON"))
			return
		}
	}
}
