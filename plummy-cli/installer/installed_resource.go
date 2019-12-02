package installer

import (
	"encoding/json"
	"fmt"
	"github.com/Masterminds/semver/v3"
)

type Resource struct {
	path string
	version *semver.Version
}

type resourceJson struct {
	Path string `json:"path"`
	Version string `json:"version"`
}

func (r *Resource) Path() string {
	return r.path
}

func (r *Resource) Version() *semver.Version {
	return r.version
}

func (r *Resource) Equals(other *Resource) bool {
	return r.path == other.Path() && r.version.Equal(other.Version())
}

func (r *Resource) MarshalJSON() ([]byte, error) {
	return json.Marshal(&resourceJson{
		Path:    r.path,
		Version: r.version.String(),
	})
}

func (r *Resource) UnmarshalJSON(data []byte) error {
	var v resourceJson
	var err error
	if err = json.Unmarshal(data, &v); err != nil {
		return err
	}
	r.path = v.Path
	if r.version, err = semver.StrictNewVersion(v.Version); err != nil {
		return fmt.Errorf("bad resource version string in JSON: %w", err)
	}
	return err
}

