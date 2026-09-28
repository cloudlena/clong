package clong

import (
	"errors"
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
		return errors.New("user ID must have between 1 and 36 characters")
	}
	if u.Name == "" || utf8.RuneCountInString(u.Name) > maxUserNameLength {
		return errors.New("username must have between 1 and 30 characters")
	}
	return nil
}
