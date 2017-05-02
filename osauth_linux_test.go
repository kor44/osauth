package osauth

import (
	"bytes"
	"os/exec"
	"testing"
)

func addUser(userName string, userPassword string, t *testing.T) {
	// create user
	cmd := exec.Command("useradd", userName, "--password", userPassword)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Skipf("Unable to create test user \"%s\":\n%s", userName, stderr.String())
	}
}

func deleteUser(userName string) {
	cmd := exec.Command("userdel", userName, "--remove")
	cmd.Run()
}
