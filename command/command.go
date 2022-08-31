package command

import (
	"fmt"
	"os"

	"github.com/dpouris/grip/locator"
	"github.com/fatih/color"
)

type Arguments struct {
	SearchString string
	Directory    string
	Options      locator.OptionConfig
}

// ParseArgs uses os.Args and parses the flags provided. If the flags(n) given are 0=n<2 then we print the usage to stdout, else we parse, assert and return the args
func ParseArgs() (Arguments, bool) {
	if len(os.Args) <= 2 {
		Usage()
		os.Exit(0)
	}

	args := os.Args[1:]
	opts := args[2:]

	mainArgs := Arguments{}

	if len(opts) != 0 {
		mainArgs.Options = parseOpts(opts)
	}

	if args[0] == "usage" {
		Usage()
		os.Exit(0)
		return Arguments{}, false
	}

	if len(args) < 2 {
		Usage()
		return Arguments{}, false
	}

	mainArgs.SearchString, mainArgs.Directory = args[0], args[1]

	ok := validateDir(mainArgs.Directory)

	if !ok {
		fmt.Printf("%s is not a valid directory\n", mainArgs.Directory)
		return Arguments{}, false
	}

	return mainArgs, true

}

// Print the basic usage on stdout
func Usage() {
	fmt.Println(
		`
Usage:

grip <searchString> ( <searchDir> | . ) [-opt]
	
Arguments:

	<searchString>	  The desired text you want to search for

	<searchDir>   	  The directory in which you'd like to search. Use '.' to search in the current directory

Options:
	
	-h 			  	  Search hidden folders and files

	`,
	)

}

// Validates that a given directory, dir, is valid
func validateDir(dir string) bool {
	_, err := os.ReadDir(dir)

	return err == nil
}

// Parses option/flag arguments given after path e.g -v, -h.
func parseOpts(opts []string) locator.OptionConfig {
	var options locator.OptionConfig
	for _, opt := range opts {
		switch opt {
		case "-h":
			options.Hidden = true
		case "-v":
			options.Verbose = true
		default:
			color.Red("Unrecognized argument %s\n", opt)
		}
	}

	return options
}
