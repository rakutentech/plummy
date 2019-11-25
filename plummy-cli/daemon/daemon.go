package daemon

import (
	"github.com/rakutentech/plummy/plummy-cli/client"
	"os"
	"runtime"
	"syscall"
	"time"
)

// Daemon is the interface for a daemon
type Daemon interface {
	Client() client.Client
	IsAlive() bool
	Stop() error
}

type localDaemon struct {
	pid int
	client client.Client
}

func (d *localDaemon) Client() client.Client {
	return d.client
}

func (d *localDaemon) IsAlive() bool {
	return isPidAlive(d.pid)
}

func (d *localDaemon) Stop() error {
	p, err := os.FindProcess(d.pid)
	if err != nil {
		return err
	}
	// Release process when done
	defer func() { _ = p.Release() }()

	// Windows does not support signals - just kill the process
	if runtime.GOOS == "windows" {
		return p.Kill()
	}

	err = p.Signal(syscall.SIGTERM)
	if err != nil {
		// Kill forcefully if we can't send a signal
		return p.Kill()
	}

	// Wait for process to shut down gracefully
	return waitAndKill(p, 10 * time.Second)
}

func waitAndKill(p *os.Process, timeout time.Duration) error {
	pid := p.Pid
	waitUntil := time.Now().Add(timeout)
	for time.Now().Before(waitUntil) {
		if !isPidAlive(pid) {
			return nil
		}
		time.Sleep(2 * time.Millisecond)
	}
	return p.Kill()
}
