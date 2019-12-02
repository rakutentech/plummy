package main

import (
	"context"
	"github.com/rakutentech/plummy/plummy-cli/client"
	"github.com/rakutentech/plummy/plummy-cli/daemon"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	options := parseOptions()
	if options.Stop {
		daemon.Stop()
		return
	}
	params, err := options.prepareParams()
	if err != nil {
		log.Fatalf("[Config Error] Cannot serialize params: %v", err)
	}

	inputs, err := options.prepareInputs()
	if err != nil {
		log.Fatalf("[Input Errror] %v", err)
	}

	d := daemon.Ensure(&daemon.StartupArgs{})
	resp, err := d.Client().Render(context.Background(), "plantuml", &client.RenderRequest{
		RawParams: params,
		Files: inputs,
	})
	if err != nil {
		log.Fatalf("[Render Errror] %v", err)
	}

	for i, file := range resp.Files {
		var err error
		if i == 0 && options.OutputFile == "-" {
			// First file is the main output - we may need to write it to stdout
			_, err = os.Stdout.Write(file.Data)
		} else {
			err = ioutil.WriteFile(file.Name, file.Data, 0777)
		}
		if err != nil {
			log.Fatalf("[Output Error] %v", err)
		}
	}
}
