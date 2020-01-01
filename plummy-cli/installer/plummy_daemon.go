package installer

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/rakutentech/plummy/plummy-cli/config"
	"log"
	"os"
	"path"
)

func FindPlummyDaemon(version *semver.Version) *Resource {
	jarFile := path.Join(jarDir(), "plummy-daemon-"+version.String()+".jar")
	if _, err := os.Stat(jarFile); err != nil {
		return nil
	}
	return &Resource{
		path:    jarFile,
		version: version,
	}
}

func UsePlummyDaemon(jarFile string) (*Resource, error) {
	if fi, err := os.Stat(jarFile); err != nil || fi.IsDir() {
		return nil, fmt.Errorf("daemon file not found at %s", jarFile)
	}
	version, err := plummySemVer(jarFile)
	if err != nil {
		return nil, err
	}
	return &Resource{
		path:    jarFile,
		version: version,
	}, nil
}

func EnsurePlummyDaemon(versionStr string) (*Resource, error) {
	version, err := parseOptionalVersion(versionStr)
	if err != nil {
		return nil, fmt.Errorf("bad pluumy version '%s' format: %w", versionStr, err)
	}
	if version.Prerelease() == "dev" {
		// We cannot find a jar file for development version versions
		return nil, nil
	}
	if d := FindPlummyDaemon(version); d != nil {
		return d, nil
	}
	return InstallPlummyDaemon(version)
}

func InstallPlummyDaemon(version *semver.Version) (*Resource, error) {
	url := fmt.Sprintf(
		"https://github.com/rakutentech/plummy/releases/download/v%v/plummy-daemon-%v.jar",
		version, version,
	)
	log.Printf("Downloading plummy daemon from %s...\n", url)
	return installFromURL(url, jarDir(), plummySemVer)
}

func plummySemVer(filename string) (*semver.Version, error) {
	ver, err := semVerFromFileName(filename, "plummy-daemon-", ".jar")
	if err != nil {
		return nil, fmt.Errorf("bad plummy daemon filename: %w", err)
	}
	return ver, nil
}

func jarDir() string {
	return path.Join(config.CacheDir(), "jar")
}

func parseOptionalVersion(versionStr string) (*semver.Version, error) {
	if versionStr == "" {
		return nil, nil // No version specified
	}
	return semver.StrictNewVersion(versionStr)
}
