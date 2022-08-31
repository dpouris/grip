package locator

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

type Locator struct {
	BaseDir string
	Options OptionConfig
}

type OptionConfig struct {
	Verbose, Hidden bool
}

type AnalyzedFile struct {
	FilePath  string
	Content   *os.File
	Locations []Location
	Ok        bool
}

type Location struct {
	LineNo   int
	Contents string
}

// Print to stdout the info of the AnalyzedFile
// e.x
//
//	fileName
//	lineNumber: foundText
func (f *AnalyzedFile) GetInfo() {
	filePath := color.YellowString(f.FilePath)
	fmt.Printf("\n%s\n", filePath)
	for _, loc := range f.Locations {
		lineNo := color.GreenString(fmt.Sprintf("%v", loc.LineNo))
		fmt.Printf("%s:%s\n", lineNo, loc.Contents)
	}
}
