package main

import (
	"github.com/alecthomas/kingpin"
	"os"
)

type options struct {
	Stop         bool
	StrVariables map[string]string
	IntVariables map[string]string
	SkinParams   map[string]string
	IncludeFiles []string
	OutputFormat string
	InputFile    string
	OutputFile   string
	ConfigLines  []string
}

func parseOptions() *options {
	app := kingpin.New("plantuml", "Plummy PlantUML Daemon Stub")
	// Support -h for help
	app.HelpFlag.Short('h')
	stop := app.Flag("stop", "Stop the daemon").Short('s').Bool()

	configLines := app.Flag("config", "Add a config Line").Short('c').Strings()
	strVariables := app.Flag("var", "Define a string variable").Short('d').StringMap()
	intVariables := app.Flag("int-var", "Define an int variable").StringMap()
	skinParams := app.Flag("skin-param", "Define a skin parameter").Short('S').StringMap()

	includeFiles := app.Flag("include",
		"Include file as if '!include file' was used, glob patterns may be used").Short('I').Strings()
	outputFormat := app.Flag("format", "Output format").Short('t').
		Default("png").
		Enum("braille", "eps", "eps:text", "latex", "latex:nopreamble", "pdf", "png", "scxml", "svg",
			"txt", "utxt", "xmi", "xmi:argo", "xmi:start")
	outputFile := app.Flag("out", "Output file ('-' for stdout)").Short('o').Default("").String()
	inputFile := app.Arg("in", "Input file").ExistingFile()

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	return &options{
		Stop:         *stop,
		ConfigLines: *configLines,
		StrVariables: *strVariables,
		IntVariables: *intVariables,
		SkinParams:   *skinParams,
		IncludeFiles: *includeFiles,
		InputFile:    *inputFile,
		OutputFormat: *outputFormat,
		OutputFile:   *outputFile,
	}
}
