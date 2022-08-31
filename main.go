package main

import (
	command "github.com/dpouris/grip/command"
	"github.com/dpouris/grip/locator"
)

func main() {
	args, ok := command.ParseArgs()

	if !ok {
		return
	}

	locator := locator.NewLocator(args.Directory)
	locator.Options = args.Options

	locator.Dig(args.SearchString)
}
