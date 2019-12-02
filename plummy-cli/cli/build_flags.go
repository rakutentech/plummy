package cli

import "fmt"

var (
	// Version specifies the built app version
	Version = "0.0.0-dev"

	// CommitHash is the git commit hash
	CommitHash = "none"

	// BuildDate is the date the app was built
	BuildDate = "unknown"
)

func VersionDescription() string {
	return fmt.Sprintf("Version %s\nGit Commit: %s\nBuild Date: %s", Version, CommitHash, BuildDate)
}