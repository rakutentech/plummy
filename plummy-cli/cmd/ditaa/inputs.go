package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/rakutentech/plummy/plummy-cli/client"
	"io/ioutil"
	"os"
)

func (opts *options) prepareInputs() ([]client.FileInput, error) {
	mainInput, err := opts.prepareMainInput()
	if err != nil {
		return nil, err
	}
	return []client.FileInput{*mainInput}, nil
}

func (opts *options) prepareMainInput() (*client.FileInput, error) {
	// Use stdin if no input file is specified
	filename := opts.InputFile
	if filename == "" {
		fi, _ := os.Stdin.Stat()
		if (fi.Mode() & os.ModeCharDevice) != 0 {
			return nil, errors.New("you can only use stdin in pipe mode")
		}
		// stdin is a pipe
		if opts.OutputFile == "" {
			opts.OutputFile = "-" // Use stdout for output
		}
		return readStdin()
	}

	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("error opening %s: %w", filename, err)
	}
	contents, err := ioutil.ReadAll(f)
	_ = f.Close()
	if err != nil {
		return nil, fmt.Errorf("error reading %s: %w", filename, err)
	}
	return &client.FileInput{
		Name:   filename,
		Reader: bytes.NewReader(contents),
	}, nil
}

func readStdin() (*client.FileInput, error) {
	// Read from stdin
	contents, err := ioutil.ReadAll(os.Stdin)
	if err != nil {
		return nil, fmt.Errorf("error reading from stdin: %w", err)
	}
	return &client.FileInput{Reader: bytes.NewReader(contents)}, nil
}

