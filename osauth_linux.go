// +build linux
package osauth

import (
	"fmt"
	"syscall"
	"unsafe"
)

/*
#cgo solaris CFLAGS: -D_POSIX_PTHREAD_SEMANTICS
#cgo LDFLAGS: -lcrypt

#include <unistd.h>
#include <stdlib.h>
#include <crypt.h>
#include <shadow.h>
#include <string.h>
*/
import "C"

func authUser(username string, password string) error {
	var pwd C.struct_spwd
	var result *C.struct_spwd
	nameC := C.CString(username)
	defer C.free(unsafe.Pointer(nameC))

	buf := alloc(userBuffer)
	defer buf.free()

	err := retryWithBuffer(buf, func() syscall.Errno {
		return syscall.Errno(C.getspnam_r(nameC,
			&pwd,
			(*C.char)(buf.ptr),
			C.size_t(buf.size),
			&result))
	})
	if err != nil {
		return fmt.Errorf("user: authenticate username %s: %v", username, err)
	}
	if result == nil {
		return UnknownUserError
	}

	passwordC := C.CString(password)
	defer C.free(unsafe.Pointer(passwordC))

	if p := C.crypt(passwordC, pwd.sp_pwdp); p == nil || C.strcmp(pwd.sp_pwdp, (*C.char)(p)) != 0 {
		return WrongPassError
	}

	return nil
}

// from here to end is copy from os/user(lookup_unix.go) package
type bufferKind C.int

const (
	userBuffer = bufferKind(C._SC_GETPW_R_SIZE_MAX)
)

func (k bufferKind) initialSize() C.size_t {
	sz := C.sysconf(C.int(k))
	if sz == -1 {
		// DragonFly and FreeBSD do not have _SC_GETPW_R_SIZE_MAX.
		// Additionally, not all Linux systems have it, either. For
		// example, the musl libc returns -1.
		return 1024
	}
	if !isSizeReasonable(int64(sz)) {
		// Truncate.  If this truly isn't enough, retryWithBuffer will error on the first run.
		return maxBufferSize
	}
	return C.size_t(sz)
}

type memBuffer struct {
	ptr  unsafe.Pointer
	size C.size_t
}

func alloc(kind bufferKind) *memBuffer {
	sz := kind.initialSize()
	return &memBuffer{
		ptr:  C.malloc(sz),
		size: sz,
	}
}

func (mb *memBuffer) resize(newSize C.size_t) {
	mb.ptr = C.realloc(mb.ptr, newSize)
	mb.size = newSize
}

func (mb *memBuffer) free() {
	C.free(mb.ptr)
}

// retryWithBuffer repeatedly calls f(), increasing the size of the
// buffer each time, until f succeeds, fails with a non-ERANGE error,
// or the buffer exceeds a reasonable limit.
func retryWithBuffer(buf *memBuffer, f func() syscall.Errno) error {
	for {
		errno := f()
		if errno == 0 {
			return nil
		} else if errno != syscall.ERANGE {
			return errno
		}
		newSize := buf.size * 2
		if !isSizeReasonable(int64(newSize)) {
			return fmt.Errorf("internal buffer exceeds %d bytes", maxBufferSize)
		}
		buf.resize(newSize)
	}
}

const maxBufferSize = 1 << 20

func isSizeReasonable(sz int64) bool {
	return sz > 0 && sz <= maxBufferSize
}
