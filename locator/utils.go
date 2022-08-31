package locator

import (
	"io/fs"
	"sync"
)

func findFiles(locator *Locator, fs []fs.DirEntry) ([]string, []string) {
	var dirs []string
	var files []string

MainLoop:
	for _, file := range fs {
		// Check if a file is a directory or not and put them in the right slice
		switch file.IsDir() {
		case true:
			// Check if we have enabled the option to search through hidden folders
			if !locator.Options.Hidden && file.Name()[0] == '.' {
				continue MainLoop
			}
			dirs = append(dirs, file.Name())
		case false:
			files = append(files, file.Name())
		}
	}

	return dirs, files
}

func runAnalyze(locator *Locator, fileName, text string, wg *sync.WaitGroup) {
	file := locator.Analyze(fileName, text)
	defer wg.Done()

	if !file.Ok {
		return
	}

	file.GetInfo()
}

func runDig(locator *Locator, text string, wg *sync.WaitGroup) {
	defer wg.Done()

	locator.Dig(text)
}
