package locator

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// Locator will be supplied with a directory as a string. We'll use the directory to recusrively search for a file that contains a given string

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

// Initializes a new Locator and returns a pointer to it. Supply "dir" with the desired base directory
func NewLocator(dir string) *Locator {
	loc := Locator{BaseDir: dir}
	return &loc
}

// Dig recursively searches through a directory and it's files to find the given string "text"
func (l *Locator) Dig(text string) {
	// Find all the subdirs inside the base dir
	fs, err := os.ReadDir(l.BaseDir)

	if err != nil {
		fmt.Println(err)
	}

	dirs, files := findFiles(l, fs)

	wg := sync.WaitGroup{}

	for _, fileName := range files {
		wg.Add(1)
		go runAnalyze(l, fileName, text, &wg)
	}

	for _, dir := range dirs {
		wg.Add(1)

		loc := NewLocator(l.BaseDir + "/" + dir)

		go runDig(loc, text, &wg)
	}

	wg.Wait()
}

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

// Use Analyze to open and scan the given file, fileName, and assert whether it contains the given string, text. If it does it is accumulated to a slice of AnalyzedFile and later returned
func (l *Locator) Analyze(fileName string, text string) AnalyzedFile {
	file, err := os.Open(l.BaseDir + "/" + fileName)

	if err != nil {
		fmt.Printf("Path \"%s\" is not valid\n", fileName)
		return AnalyzedFile{}
	}

	// old way
	// var readFile string
	// buffer := make([]byte, 10)
	// for {
	// 	writtenBytes, err := file.Read(buffer)

	// 	if err == io.EOF {
	// 		break
	// 	}

	// 	readFile += string(buffer[:writtenBytes])

	// 	exists := strings.Contains(readFile, text)

	// 	if exists {
	// 		resString := strings.Replace(readFile, text, color.RedString(text), -1)
	// 		return &AnalyzedFile{FileName: fileName, Content: file, Read: resString, Ok: true}
	// 	}

	// }

	scanner := bufio.NewScanner(file)
	lineNo := 1
	var locations []Location

	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		line := scanner.Text()

		exists := strings.Contains(line, text)

		if exists {
			resString := strings.Replace(line, text, color.RedString(text), -1)
			loc := Location{LineNo: lineNo, Contents: resString}
			locations = append(locations, loc)
		}

		lineNo++
	}

	analyzedFile := AnalyzedFile{FilePath: l.BaseDir + "/" + fileName, Content: file, Locations: locations}

	if len(analyzedFile.Locations) > 0 {
		analyzedFile.Ok = true
	}

	return analyzedFile
}
