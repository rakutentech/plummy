package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
)

type globalParams struct {
	Format      string            `json:"format,omitempty"`
	ConfigLines []string          `json:"configLines,omitempty"`
}

func (opts *options) prepareParams() ([]byte, error) {
	configLines := append([]string{}, opts.ConfigLines...)

	// Add configuration lines based on skin params and includes
	for name, value := range opts.StrVariables {
		configLines = append(configLines, fmt.Sprintf("!$%s = \"%s\"", name, value))
	}
	for name, value := range opts.IntVariables {
		configLines = append(configLines, fmt.Sprintf("!$%s = %s", name, value))
	}

	configLines, err := opts.addIncludes(configLines)
	if err != nil {
		return nil, err
	}

	for name, value := range opts.SkinParams {
		configLines = append(configLines, fmt.Sprintf("skinparamlocked %s %s", name, value))
	}

	params := globalParams{
		Format:      opts.OutputFormat,
		ConfigLines: configLines,
	}
	return json.Marshal(&params)
}

func (opts *options) addIncludes(configLines []string) ([]string, error) {
	for _, includePattern := range opts.IncludeFiles {
		matches, err := filepath.Glob(includePattern)
		if err != nil {
			// If glob fails will just try to use the pattern directly as a filename.
			matches = []string{includePattern}
		}
		for _, filename := range matches {
			contents, err := ioutil.ReadFile(filename)
			if err != nil {
				return nil, fmt.Errorf("cannot include file %s: %w", filename, err)
			}
			configLines = append(configLines, string(contents))
		}
	}
	return configLines, nil
}
