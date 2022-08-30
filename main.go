package main

import (
	command "github.com/dpouris/grip/command"
	"github.com/dpouris/grip/locator"
)

func main() {
	var args []string
	var ok bool

	args, ok = command.ParseArgs()

	if !ok {
		return
	}

	searchString := args[0]
	searchDir := args[1]

	locator := locator.NewLocator(searchDir)

	locator.Dig(searchString)
}
