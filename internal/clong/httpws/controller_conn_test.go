package httpws_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
	"github.com/cloudlena/clong/internal/clong/httpws"
)

func TestHandleControllerConnRejectsBadCookies(t *testing.T) {
	tests := map[string]struct {
		cookies []*http.Cookie
		want    int
	}{
		"no cookies":           {nil, http.StatusUnauthorized},
		"missing username":     {[]*http.Cookie{{Name: "userid", Value: "abc123"}}, http.StatusUnauthorized},
		"missing user ID":      {[]*http.Cookie{{Name: "username", Value: "Alice"}}, http.StatusUnauthorized},
		"undecodable username": {[]*http.Cookie{{Name: "userid", Value: "abc123"}, {Name: "username", Value: "Alice%zz"}}, http.StatusBadRequest},
		"username too long":    {[]*http.Cookie{{Name: "userid", Value: "abc123"}, {Name: "username", Value: strings.Repeat("a", 31)}}, http.StatusBadRequest},
		"user ID too long":     {[]*http.Cookie{{Name: "userid", Value: strings.Repeat("a", 37)}, {Name: "username", Value: "Alice"}}, http.StatusBadRequest},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/ws/controller", nil)
			for _, c := range tt.cookies {
				req.AddCookie(c)
			}
			w := httptest.NewRecorder()
			httpws.HandleControllerConn(clong.NewService(&mockScoreStore{}))(w, req)

			if w.Code != tt.want {
				t.Errorf("expected %d, got %d", tt.want, w.Code)
			}
		})
	}
}
