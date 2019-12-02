package installer

import (
	"github.com/Masterminds/semver/v3"
	"github.com/rakutentech/plummy/plummy-cli/config"
	"log"
	"path"
	"path/filepath"
)

func PlummyDaemon() *Resource {
	version, filename := findPlummyDaemon()
	if version == nil || filename == "" {
		return nil
	}
	return &Resource{
		path:    filename,
		version: version,
	}
}

func EnsurePlummyDaemon() (*Resource, error) {
	if daemon := PlummyDaemon(); daemon != nil {
		return daemon, nil
	}
	return InstallPlummyDaemon()
}

func InstallPlummyDaemon() (*Resource, error) {
	url := "https://github.com/rakutentech/plummy/releases/download/v0.1.0/plummy-daemon-0.1.0.jar"
	log.Printf("Downloading plummy daemon from %s...\n", url)
	return installFromURL(url, jarDir(), plummySemVer)
}

func findPlummyDaemon() (*semver.Version, string) {
	jarPattern := path.Join(jarDir(), "plummy-daemon-*.jar")
	matches, err := filepath.Glob(jarPattern)
	if err != nil {
		return nil, ""
	}
	var latestVersion *semver.Version
	var latestVersionFile string
	for _, filename := range matches {
		version, err := plummySemVer(filename)
		if err != nil {
			continue
		}
		if latestVersion == nil || version.GreaterThan(latestVersion) {
			latestVersion = version
			latestVersionFile = filename
		}
	}
	return latestVersion, latestVersionFile
}

func plummySemVer(filename string) (*semver.Version, error) {
	return semVerFromFileName(filename, "plummy-daemon-", ".jar")
}

func jarDir() string {
	return path.Join(config.CacheDir(), "jar")
}
