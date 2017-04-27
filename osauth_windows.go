// +build win32
package osauth

import (
	"syscall"
	"unsafe"
)

var (
	advapi32       = syscall.NewLazyDLL("advapi32.dll")
	procLogonUserW = advapi32.NewProc("LogonUserW")
)

const (
	logon32LogonNetwork    = 3
	logon32ProviderDefault = 0
)

func authUser(username string, password string) error {
	var token syscall.Handle
	r1, _, err := procLogonUserW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(""))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(password))),
		uintptr(logon32LogonNetwork),
		uintptr(logon32ProviderDefault),
		uintptr(unsafe.Pointer(&token)))

	if int(r1) == 0 {
		return err
	}

	return nil
}
