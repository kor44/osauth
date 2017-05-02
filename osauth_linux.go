// +build linux
package osauth

import (
	"errors"
	"unsafe"
)

/*
#cgo LDFLAGS: -lcrypt

#define _GNU_SOURCE

#include <stdlib.h>
#include <string.h>
#include <crypt.h>
#include <shadow.h>

int check_auth(char *username, char *password) {
	struct spwd *usersp;
	// This call is needed to get the password; otherwise 'x' is returned
    if ((usersp= getspnam(username)) == NULL)
    {
        //snprintf(errbuf, PCAP_ERRBUF_SIZE, "Authentication failed: no such user");
        return 2;
    }

    if (strcmp(usersp->sp_pwdp, (char *) crypt(password, usersp->sp_pwdp) ) != 0)
    {
		//printf("usersp->sp_pwdp=%s\n", usersp->sp_pwdp);
		//printf("crypt =%s\n", (char *) crypt(password, usersp->sp_pwdp));
        //snprintf(errbuf, PCAP_ERRBUF_SIZE, "Authentication failed: password incorrect");
        return 0;
    }

	return 1;
}
*/
import "C"

func authUser(username string, password string) error {
	c_name := C.CString(username)
	defer C.free(unsafe.Pointer(c_name))

	c_pass := C.CString(password)
	defer C.free(unsafe.Pointer(c_pass))

	result := C.int(C.check_auth(c_name, c_pass))
	if result == 1 {
		return nil
	}

	return errors.New("Authentication failed: password incorrect")
}
