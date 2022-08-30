package command

import (
	"fmt"
	"os"
)

// flag parser
func ParseArgs() ([]string, bool) {
	args := os.Args[1:]

	if args[0] == "usage" {
		Usage()
		os.Exit(0)
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

func Usage() {
	fmt.Println(
		`
		
  Usage:
  
 	mgrep <searchString> ( <searchDir> | . )

	- searchString	  The desired text you want to search for

	- searchDir   	  The directory in which you'd like to search for the specified text
	
	Alternatevely you can use [ . ] to search in the current directory

	`,
	)

}

func validateDir(dir string) bool {
	_, err := os.ReadDir(dir)

	return err == nil
}
