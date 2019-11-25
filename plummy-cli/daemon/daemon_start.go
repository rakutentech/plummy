package daemon

import (
	"context"
	"github.com/rakutentech/plummy/plummy-cli/jvm"
	"log"
	"os"
	"path"
	"time"
)

// Find gets the existing daemon or returns nil if no active daemon is found.
func Find() Daemon {
	spec, err := readDaemonFile()
	if err != nil {
		log.Printf("[WARN] Can't read daemon spec file: %v\n", err)
		return nil
	}

	// May still return nil if no daemon file is found or daemon is stopped, but without warning
	return spec.ToDaemon()
}

// Ensure returns an active daemon or starts a new one if
func Ensure() Daemon {
	d := Find()

	if d == nil || !d.IsAlive() {
		d = start()
	}

	waitForDaemon(d)

	return d
}

func Stop() {
	daemon := Find()
	if daemon != nil && daemon.IsAlive() {
		log.Printf("Stopping daemon")
		err := daemon.Stop()
		if err != nil {
			log.Printf("[ERROR] Cannot stop daemon: %s", err.Error())
		}
	} else {
		log.Printf("Daemon not found")
	}
}

func waitForDaemon(d Daemon) {
	timeoutCtx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	if err := d.Client().WaitReady(timeoutCtx); err != nil {
		log.Fatalf("Daemon not ready: %v\n", err)
	}
}

func start() Daemon {
	// TODO: Properly install file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Home directory not found: %v\n", err)
	}
	daemonJar := path.Join(homeDir, "plummy-0.1.0.jar")
	if !pathExists(daemonJar) {
		log.Fatalf("Daemon file %s not found\n", daemonJar)
	}
	java := jvm.Default()
	if java == nil {
		log.Fatal("No JVM installation detected - please set JAVA_HOME\n")
	}
	proc, err := java.Daemonize(jvm.DaemonOptions{}, "-jar", daemonJar) // TODO: Allow to Port
	if err != nil {
		log.Fatalf("Could not start daemon: %v\n", err)
	}
	spec := &daemonSpec{
		Pid:     proc.Pid,
		BaseURL: "http://localhost:4545/",
	}
	if err := writeDaemonFile(spec); err != nil {
		log.Printf("[WARN] Can't write daemon spec file: %v\n", err)
	}

	return spec.ToDaemon()
}
