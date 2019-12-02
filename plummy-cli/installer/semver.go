package installer

import (
	"fmt"
	"github.com/Masterminds/semver/v3"
	"path"
	"strings"
)

func semVerFromFileName(filename, prefix, suffix string) (*semver.Version, error) {
	filename = path.Base(filename)
	if !strings.HasPrefix(filename, prefix) {
		return nil, fmt.Errorf("expected filename '%s' to have prefix '%s'", filename, prefix)
	}
	if !strings.HasSuffix(filename, suffix) {
		return nil, fmt.Errorf("expected filename '%s' to have suffix '%s'", filename, suffix)
	}
	versionStr := filename[len(prefix) : len(filename)-len(suffix)]
	return semver.StrictNewVersion(versionStr)
}

