package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/rakutentech/plummy/plummy-cli/client"
	"github.com/rakutentech/plummy/plummy-cli/config"
	"github.com/rakutentech/plummy/plummy-cli/installer"
	"io/ioutil"
	"os"
	"path"
)

type daemonSpec struct {
	Pid     int                 `json:"pid"`
	BaseURL string              `json:"base_url"`
	Jar     *installer.Resource `json:"resource"`
}

func (spec *daemonSpec) ToDaemon() Daemon {
	if spec.Pid == 0 {
		return nil
	}
	return &localDaemon{
		pid:    spec.Pid,
		client: client.NewHttpClient(spec.BaseURL),
		version: spec.Jar.Version(),
	}
}

func readDaemonFile() (*daemonSpec, error) {
	result := &daemonSpec{}
	filename, err := daemonFilename()
	if err != nil {
		return nil, err
	}
	if !pathExists(filename) {
		return result, nil // File not found return an empty spec
	}
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func writeDaemonFile(spec *daemonSpec) error {
	filename, err := daemonFilename()
	if err != nil {
		return err
	}
	specBytes, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, specBytes, 0644)
	if err != nil {
		return err
	}
	return nil
}

func daemonFilename() (string, error) {
	// Ensure directory exists
	cacheDir := config.CacheDir()
	if err := config.EnsureDir(cacheDir); err != nil {
		return "", fmt.Errorf("cannot create cache directory: %w", err)
	}
	return path.Join(cacheDir, "daemon.json"), nil
}

func pathExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}
