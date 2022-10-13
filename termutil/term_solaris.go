//go:build solaris && !appengine
// +build solaris,!appengine

package termutil

const (
	ioctlReadTermios  = 0x5401 // syscall.TCGETS
	ioctlWriteTermios = 0x5402 // syscall.TCSETS
	sysIoctl          = 54
)
