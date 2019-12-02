package installer

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"github.com/rakutentech/plummy/plummy-cli/config"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
)

func installFromURL(remoteURL, targetDir string, semVerFunc func(filename string) (*semver.Version, error)) (*Resource, error) {
	parsedURL, err := url.Parse(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("cannot parse url: %w", err)
	}

	resp, err := http.Get(remoteURL)
	if err != nil {
		return nil, fmt.Errorf("cannot download %s: %w", remoteURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("resource at %s returned status %d %s", remoteURL, resp.StatusCode, resp.Status)
	}

	// Determine filename, path, version
	targetFilename := path.Base(parsedURL.Path)
	targetPath := path.Join(targetDir, targetFilename)
	version, err := semVerFunc(targetFilename)
	if err != nil {
		return nil, err
	}

	// Ensure target directory exists
	if err := config.EnsureDir(targetDir); err != nil {
		return nil, fmt.Errorf("cannot create target directory: %w", err)
	}

	// Create target file
	out, err := os.Create(targetPath)
	if err != nil {
		return nil, err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return nil, err
	}

	return &Resource{
		path:    targetPath,
		version: version,
	}, nil
}
