//go:build (darwin || freebsd || netbsd || openbsd || dragonfly) && !appengine
// +build darwin freebsd netbsd openbsd dragonfly
// +build !appengine

package termutil

import "syscall"

const (
	ioctlReadTermios  = syscall.TIOCGETA
	ioctlWriteTermios = syscall.TIOCSETA
)
