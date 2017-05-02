package osauth

import (
	"testing"
)

func TestAuth(t *testing.T) {
	var userName = "osauth_testuser"
	var userPassword = "osauth_pass"

	addUser(userName, userPassword, t)
	defer deleteUser(userName)

	cases := map[string]struct {
		username string
		pass     string
		result   func(error) bool
	}{
		"Auth Successful": {userName, userPassword, func(err error) bool { return err == nil }},
		"Auth failed":     {userName, "wrong_password", func(err error) bool { return err != nil && err != UnknownUserError }},
	}

	for k, tc := range cases {
		t.Run(k, func(t *testing.T) {
			if err := AuthUser(tc.username, tc.pass); !tc.result(err) {
				t.Error(err)
				t.Fail()
			}
		})

	}
}

func TestAuthUserNotExist(t *testing.T) {
	if err := AuthUser("wrong_user", "wrong_password"); err != UnknownUserError {
		t.Error(err)
		t.Fail()
	}
}
