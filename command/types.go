package command

import "github.com/dpouris/grip/locator"

type Arguments struct {
	SearchString string
	Directory    string
	Options      locator.OptionConfig
}

type Parser interface {
	ParseArgs() (Arguments, bool)
	ParseOpts(opts []string) locator.OptionConfig
}
