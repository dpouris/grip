package command

import (
	"fmt"
	"os"
)

type Arguments struct {
	SearchString string
	Directory    string
	Options      map[string]bool
}

// ParseArgs uses os.Args and parses the flags provided. If the flags(n) given are 0=n<2 then we print the usage to stdout, else we parse, assert and return the args
func ParseArgs() (Arguments, bool) {
	if len(os.Args) <= 1 {
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
  
 	grip <searchString> ( <searchDir> | . )

	- searchString	  The desired text you want to search for

	- searchDir   	  The directory in which you'd like to search for the specified text
	
	Alternatevely you can use [ . ] to search in the current directory

	`,
	)

}

// Validates that a given directory, dir, is valid
func validateDir(dir string) bool {
	_, err := os.ReadDir(dir)

	return err == nil
}

// Parses options/flags arguments given after path e.g -v, -h.
func parseOpts(opts []string) map[string]bool {
	options := map[string]bool{"-h": false, "-a": false}
	for _, opt := range opts {
		switch opt {
		case "-h":
			options[opt] = true
		case "-a":
			options[opt] = true
		default:
			fmt.Printf("Unrecognized argument %s", opt)
		}
	}

	return options
}
