package jvm

import (
	"bytes"
	"io/ioutil"
	"os/exec"
)

func captureStdout(command string, args... string) (*bytes.Buffer, error) {
	stdout := &bytes.Buffer{}
	cmd := exec.Command(command, args...)
	cmd.Stdout = stdout
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return stdout, nil
}

func captureStdoutString(command string, args... string) (string, error) {
	stdout, err := captureStdout(command, args...)
	if stdout == nil {
		return "", err
	}
	return stdout.String(), err
}

func captureStderr(command string, args... string) (*bytes.Buffer, error) {
	stderr := &bytes.Buffer{}
	cmd := exec.Command(command, args...)
	cmd.Stderr = stderr
	cmd.Stdout = ioutil.Discard
	err := cmd.Run()
	if err != nil {
		return nil, err
	}
	return stderr, nil
}

