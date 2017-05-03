package osauth

import (
	"errors"
	"os/user"
)

var (
	UnknownUserError = errors.New("Authentication failed: no such user")
	WrongPassError   = errors.New("Authentication failed: password incorrect")
)

// Check authentication. If user name not found return UnknownUserError,
// if password is incorrect return WrongPassError, or platform specific error
func AuthUser(username string, passsword string) error {
	if _, err := user.Lookup(username); err != nil {
		return UnknownUserError
	}

	return authUser(username, passsword)
}
