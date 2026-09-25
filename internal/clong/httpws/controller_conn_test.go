package httpws_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
)

func TestHandleControllerConnMissingCookies(t *testing.T) {
	tests := map[string][]*http.Cookie{
		"no cookies":       nil,
		"missing username": {{Name: "userid", Value: "abc123"}},
		"missing user ID":  {{Name: "username", Value: "Alice"}},
	}

	for name, cookies := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ws/controller", nil)
			for _, c := range cookies {
				req.AddCookie(c)
			}
			w := httptest.NewRecorder()
			httpws.HandleControllerConn(clong.NewService(&mockScoreStore{}))(w, req)

			if w.Code != http.StatusUnauthorized {
				t.Errorf("expected 401, got %d", w.Code)
			}
		})
	}
}
