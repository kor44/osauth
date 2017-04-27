package osauth

import (
	"bytes"
	"os/exec"
	"testing"
)

func addUser(userName string, userPassword string, t *testing.T) {
	// set active console code page to United States
	exec.Command("chcp", "437").Run()

	// create user
	cmd := exec.Command("net", "user", userName, userPassword, "/ADD")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Skipf("Unable to create test user \"%s\":\n%s", userName, stderr.String())
	}
}

func deleteUser(userName string) {
	cmd := exec.Command("net", "user", userName, "/DELETE")
	cmd.Run()
}

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
