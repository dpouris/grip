package main

import (
	command "mgrep/cli"
	"mgrep/locator"
	"os"
)

func main() {
	var args []string
	var ok bool

	args, ok = command.ParseArgs()

	if !ok {
		os.Exit(1)
	}

	searchString := args[0]
	searchDir := args[1]

	locator := locator.NewLocator(searchDir)

	locator.Dig(searchString)
}
