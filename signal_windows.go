//go:build windows
// +build windows

package nssh

import "syscall"

const SIGWINCH = syscall.Signal(0xff)
