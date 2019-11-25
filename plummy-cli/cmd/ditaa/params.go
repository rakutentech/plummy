package main

import (
	"encoding/json"
)

type globalParams struct {
	Debug      bool             `json:"debug"`
	Rendering  renderingParams  `json:"rendering"`
	Processing processingParams `json:"processing"`
}

type renderingParams struct {
	BackgroundColor string `json:"backgroundColor,omitempty"`
	Format string `json:"format,omitempty"`
}

type processingParams struct {
	Verbose bool `json:"verbose"`
}

func (opts *options) prepareParams() ([]byte, error) {
	params := globalParams{
		Debug: opts.Debug,
		Rendering: renderingParams{
			BackgroundColor: opts.Background,
			Format: opts.OutputFormat,
		},
		Processing: processingParams{
			Verbose: opts.Verbose,
		},
	}
	return json.Marshal(&params)
}
