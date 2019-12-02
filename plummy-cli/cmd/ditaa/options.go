package main

import (
	"github.com/alecthomas/kingpin"
	"github.com/rakutentech/plummy/plummy-cli/cli"
	"os"
)

type options struct {
	Stop         bool
	Debug        bool
	Verbose      bool
	OutputFormat string
	InputFile    string
	OutputFile   string
	Background   string
}

func parseOptions() *options {
	app := kingpin.New("ditaa", "Plummy Ditaa Daemon Stub")
	app.Version(cli.VersionDescription())
	// Support -h for help
	app.HelpFlag.Short('h')
	stop := app.Flag("stop", "Stop the daemon").Short('s').Bool()
	debug := app.Flag("debug", "Renders the debug grid over the resulting image").Short('d').Bool()
	verbose := app.Flag("verbose", "Makes ditaa more verbose").Short('v').Bool()
	background := app.Flag("background", "The background colour of the image. The format "+
		"should be a six-digit hexadecimal number (as in HTML, FF0000 for red). Pass an eight-digit "+
		"hex to define transparency. This is overridden by --transparent.").Short('b').String()
	transparent := app.Flag("transparent", "Causes the diagram to be rendered on a "+
		"transparent background. Overrides").Short('T').Bool()

	svg := app.Flag("svg", "Use the SVG output format").Bool()
	outputFormat := app.Flag("format", "Choose output format").Short('t').
		Default("png").
		Enum("png", "svg")
	outputFile := app.Flag("out", "Output file ('-' for stdout)").Short('o').Default("").String()
	inputFile := app.Arg("in", "Input file").ExistingFile()

	// Support the --svg flag if present
	if *svg {
		*outputFile = "svg"
	}

	if *transparent {
		*background = "00000000" // Black with 0 Alpha
	}

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	return &options{
		Stop:         *stop,
		Debug:        *debug,
		Verbose:      *verbose,
		InputFile:    *inputFile,
		OutputFormat: *outputFormat,
		OutputFile:   *outputFile,
		Background:   *background,
	}
}
