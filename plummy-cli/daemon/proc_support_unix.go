// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris plan9

package daemon

import "syscall"

func isPidAlive(pid int) bool {
	// Send a 0 signal to the process to check if it's alive.
	return syscall.Kill(pid, 0) == nil
}
