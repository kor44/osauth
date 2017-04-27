package osauth

import (
	"errors"
	"os/user"
)

var UnknownUserError = errors.New("Authentication failed: no such user")

// Check authentication. If user name not found return UnknownUserError
func AuthUser(username string, passsword string) error {
	if _, err := user.Lookup(username); err != nil {
		return UnknownUserError
	}

	return authUser(username, passsword)
}
