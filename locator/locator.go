package locator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// Locator will be supplied with a directory as a string. We'll use the directory to recusrively search for a file that contains a given string

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
