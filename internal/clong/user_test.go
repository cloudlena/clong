package clong_test

import (
	"strings"
	"testing"

	"github.com/cloudlena/clong/internal/clong"
)

func TestUserValidate(t *testing.T) {
	tests := map[string]struct {
		user  clong.User
		valid bool
	}{
		"valid":          {clong.User{ID: "abc", Name: "Alice"}, true},
		"multibyte name": {clong.User{ID: "abc", Name: strings.Repeat("ö", 30)}, true},
		"missing ID":     {clong.User{Name: "Alice"}, false},
		"ID too long":    {clong.User{ID: strings.Repeat("a", 37), Name: "Alice"}, false},
		"missing name":   {clong.User{ID: "abc"}, false},
		"name too long":  {clong.User{ID: "abc", Name: strings.Repeat("a", 31)}, false},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			err := tt.user.Validate()
			if tt.valid && err != nil {
				t.Errorf("expected valid user, got %v", err)
			}
			if !tt.valid && err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}
