package httpws

import (
	"net/http"
	"testing"
)

func TestCookieValFound(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "userid", Value: "abc123"},
		{Name: "username", Value: "Alice"},
	}

	val, ok := cookieVal(cookies, "userid")
	if !ok {
		t.Fatal("expected cookie to be found")
	}
	if val != "abc123" {
		t.Errorf("expected 'abc123', got %q", val)
	}
}

func TestCookieValNotFound(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "userid", Value: "abc123"},
	}

	_, ok := cookieVal(cookies, "username")
	if ok {
		t.Error("expected cookie not to be found")
	}
}

func TestCookieValEmptySlice(t *testing.T) {
	_, ok := cookieVal([]*http.Cookie{}, "userid")
	if ok {
		t.Error("expected false for empty cookie slice")
	}
}

func TestCookieValReturnsFirstMatch(t *testing.T) {
	cookies := []*http.Cookie{
		{Name: "theme", Value: "dark"},
		{Name: "theme", Value: "light"},
	}

	val, ok := cookieVal(cookies, "theme")
	if !ok {
		t.Fatal("expected cookie to be found")
	}
	if val != "dark" {
		t.Errorf("expected first match 'dark', got %q", val)
	}
}
