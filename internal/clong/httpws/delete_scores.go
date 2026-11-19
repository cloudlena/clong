package httpws

import (
	"log"
	"net/http"

	"github.com/cloudlena/clong/internal/clong"
)

// HandleDeleteScores deletes all scores and resets the scoreboards.
func HandleDeleteScores(svc *clong.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := svc.ResetScores(r.Context()); err != nil {
			log.Println(err)
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}
}
