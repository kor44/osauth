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
