package daemon

import (
	"encoding/json"
	"fmt"
	"github.com/rakutentech/plummy/plummy-cli/client"
	"io/ioutil"
	"os"
	"path"
)

type daemonSpec struct {
	Pid     int    `json:"pid"`
	BaseURL string `json:"base_url"`
}

func (spec *daemonSpec) ToDaemon() Daemon {
	if spec.Pid == 0 {
		return nil
	}
	return &localDaemon{
		pid: spec.Pid,
		client: client.NewHttpClient(spec.BaseURL),
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
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return "", fmt.Errorf("cannot get user cache dir: %w", err)
	}
	daemonDir := path.Join(cacheDir, "plummy")

	// Ensure directory exists
	if err := os.MkdirAll(daemonDir, 0755); err != nil {
		return "", fmt.Errorf("cannot create cache directory: %w", err)
	}
	return path.Join(daemonDir, "daemon.json"), nil
}

func pathExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}