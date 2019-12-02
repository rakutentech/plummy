package daemon

import (
	"context"
	"github.com/rakutentech/plummy/plummy-cli/cli"
	"github.com/rakutentech/plummy/plummy-cli/installer"
	"github.com/rakutentech/plummy/plummy-cli/jvm"
	"log"
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

// Ensure returns an active daemon or starts a new one if necessary
func Ensure(args *StartupArgs) Daemon {
	var jar *installer.Resource
	var err error
	if args.JarFile != "" {
		jar, err = installer.UsePlummyDaemon(args.JarFile)
		if err != nil {
			log.Fatalf("Daemon not found: %v\n", err)
		}
	} else {
		if args.Version == "" {
			args.Version = cli.Version
		}
		jar, err = installer.EnsurePlummyDaemon(args.Version)
		if err != nil {
			log.Fatalf("Daemon Installation failed: %v\n", err)
		}
	}
	d := Find()
	if d == nil || !d.IsAlive() {
		d = start(jar)
	} else if shouldRestart(d, jar) {
		// Restart the daemon if user asked for a different version a different version
		log.Printf("Current daemon version is %s - restarting daemon with version %s\n", d.Version(), jar.Version())
		d = restart(d, jar)
	}

	waitForDaemon(d)

	return d
}

func Stop() {
	daemon := Find()
	if daemon != nil && daemon.IsAlive() {
		log.Printf("Stopping daemon\n")
		err := daemon.Stop()
		if err != nil {
			log.Printf("[ERROR] Cannot stop daemon: %s\n", err.Error())
		}
	} else {
		log.Println("Daemon not found")
	}
}

func waitForDaemon(d Daemon) {
	timeoutCtx, _ := context.WithTimeout(context.Background(), 15*time.Second)
	if err := d.Client().WaitReady(timeoutCtx); err != nil {
		log.Fatalf("Daemon not ready: %v\n", err)
	}
}

func start(daemonJar *installer.Resource) Daemon {
	java := jvm.Default()
	if java == nil {
		log.Fatal("No JVM installation detected - please set JAVA_HOME\n")
	}
	proc, err := java.Daemonize(jvm.DaemonOptions{}, "-jar", daemonJar.Path()) // TODO: Allow to specify Port
	if err != nil {
		log.Fatalf("Could not start daemon: %v\n", err)
	}
	spec := &daemonSpec{
		Pid:     proc.Pid,
		BaseURL: "http://localhost:4545/",
		Jar:     daemonJar,
	}
	if err := writeDaemonFile(spec); err != nil {
		log.Printf("[WARN] Can't write daemon spec file: %v\n", err)
	}

	return spec.ToDaemon()
}

func restart(oldDaemon Daemon, newDaemonJar *installer.Resource) Daemon {
	err := oldDaemon.Stop()
	if err != nil {
		log.Printf("[WARN] Can't stop old daemon: %v\n", err)
	}
	return start(newDaemonJar)
}

func shouldRestart(oldDaemon Daemon, newDaemonJar *installer.Resource) bool {
	v := oldDaemon.Version()
	if v == nil {
		return true
	}
	return !v.Equal(newDaemonJar.Version())
}
