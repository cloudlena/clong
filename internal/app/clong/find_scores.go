package clong

import (
	"encoding/json"
	"net/http"

	"github.com/mastertinner/clong/internal/app/clong/scores"
	"github.com/pkg/errors"
)

// HandleFindScores returns all scores as JSON.
func HandleFindScores(store scores.ScoreStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve scores from DB
		ctx := r.Context()
		scrs, err := store.Scores(ctx)
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
