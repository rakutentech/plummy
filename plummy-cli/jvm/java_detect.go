package jvm

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func dirVersion(path string) (*Jvm, error) {
	javaExec := getJavaExec(path)
	if javaExec == "" {
		return nil, fmt.Errorf("bad java path '%s'", path)
	}
	version, err := ExecVersion(javaExec)
	if err != nil {
		return nil, err
	}
	return &Jvm{Binary: javaExec, Version: version}, nil
}

func Default() *Jvm {
	homeDir := os.Getenv("JAVA_HOME")
	inst, err := dirVersion(homeDir)
	if err == nil {
		return inst
	}
	pathExec, err := exec.LookPath("java")
	if err == nil {
		if version, err := ExecVersion(pathExec); err == nil {
			return &Jvm{Binary: pathExec, Version: version}
		}
	}
	if runtime.GOOS == "darwin" && isFile("/usr/libexec/java_home") {
		homeDir, _ = captureStdoutString("/usr/libexec/java_home")
		homeDir = strings.TrimSpace(homeDir)
		if homeDir != "" {
			inst, err = dirVersion(homeDir)
			if err == nil {
				return inst
			}
		}
	}
	return nil
}
