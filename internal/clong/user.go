package clong

import (
	"fmt"
	"unicode/utf8"
)

const (
	maxUserIDLength   = 36
	maxUserNameLength = 30
)

// User is a person who interacts with the app.
type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Validate checks whether a user can be stored.
func (u User) Validate() error {
	if u.ID == "" || utf8.RuneCountInString(u.ID) > maxUserIDLength {
		return fmt.Errorf("user ID must have between 1 and %d characters", maxUserIDLength)
	}
	if u.Name == "" || utf8.RuneCountInString(u.Name) > maxUserNameLength {
		return fmt.Errorf("username must have between 1 and %d characters", maxUserNameLength)
	}
	return nil
}
