package daemon

import "os"

func isPidAlive(pid int) bool {
	// FindProcess() will return an error if process is not alive on Windows
	p, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	p.Release()
	return true
}
