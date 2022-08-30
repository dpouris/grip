package command

import (
	"fmt"
	"os"
)

// ParseArgs uses os.Args and parses the flags provided. If the flags(n) given are 0=n<2 then we print the usage to stdout, else we parse, assert and return the args
func ParseArgs() ([]string, bool) {
	if len(os.Args) <= 1 {
		Usage()
		os.Exit(0)
	}

	args := os.Args[1:]

	if args[0] == "usage" {
		Usage()
		os.Exit(0)
		return nil, false
	}

	if len(args) < 2 {
		Usage()
		return nil, false
	}

	searchString, dir := args[0], args[1]

	ok := validateDir(dir)

	if !ok {
		fmt.Printf("%s is not a valid directory\n", dir)
		return nil, false
	}

	return []string{searchString, dir}, true

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
