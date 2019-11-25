package jvm

import (
	"io"
	"os"
	"os/exec"
)

type Jvm struct {
	Binary  string
	Version Version
}

type Process struct {
	*os.Process
}

type DaemonOptions struct {
	WriteStdErr bool
}

func (j *Jvm) Command(args ...string) *exec.Cmd {
	return exec.Command(j.Binary, args...)
}

func (j *Jvm) Daemonize(options DaemonOptions, args ...string) (*Process, error) {
	cmd := j.Command(args...)
	// TODO: Properly track stderr/stdout
	stderr, _ := cmd.StderrPipe()
	if options.WriteStdErr {
		go func() {
			_, _ = io.Copy(os.Stderr, stderr)
		}()
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Process{
		Process: cmd.Process,
	}, nil
}
