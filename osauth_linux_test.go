package osauth

import (
	"bytes"
	"fmt"
	"os/exec"
	"testing"
)

func addUser(userName string, userPassword string, t *testing.T) {
	// create user
	useradd := exec.Command("useradd", userName)
	var stderr bytes.Buffer
	useradd.Stderr = &stderr
	if err := useradd.Run(); err != nil {
		t.Skipf("Unable to create test user \"%s\":\n%s", userName, stderr.String())
	}

	if userPassword == "" {
		return
	}

	// set password
	chpasswd := exec.Command("bash", "-c", fmt.Sprintf(`echo "%s:%s" | chpasswd`, userName, userPassword))
	chpasswd.Stderr = &stderr
	if err := chpasswd.Run(); err != nil {
		t.Skipf("Unable to set test user password \"%s\":\n%s", userName, stderr.String())
	}

}

func deleteUser(userName string) {
	cmd := exec.Command("userdel", userName, "--remove")
	cmd.Run()
}
