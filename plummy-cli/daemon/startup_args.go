package daemon

// StartupArgs describe the arguments used for starting app a daemon
type StartupArgs struct {
	// Port is the TCP port the daemon will use for listening
	Port int

	// JarFile specifies a custom jar file to use for running the daemon
	JarFile string

	// Version specifies an explicit daemon version to download and use
	Version string
}
