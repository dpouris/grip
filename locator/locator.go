package locator

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/fatih/color"
)

// Locator will be supplied with a directory as a string. We'll use the directory to recusrively search for a file that contains a given string

type Locator struct {
	BaseDir string
}

type AnalyzedFile struct {
	FileName   string
	Content    *os.File
	Read       string
	LineNumber int
}

// Print to stdout the info of the AnalyzedFile
// e.x
//
//	fileName
//	lineNumber: foundText
func (f *AnalyzedFile) GetInfo() {
	fileName := color.YellowString(f.FileName)
	lineNo := color.GreenString(fmt.Sprintf("%v", f.LineNumber))
	fmt.Printf("\n%s\n%s:%s\n", fileName, lineNo, f.Read)
}

// Initializes a new Locator and returns a pointer to it. Supply "dir" with the desired base directory
func NewLocator(dir string) *Locator {
	loc := Locator{BaseDir: dir}
	return &loc
}

// Dig recursively searches through a directory and it's files to find the given string "text"
func (l *Locator) Dig(text string) {
	// Find all the subdirs inside the base dir
	dirFs, err := ioutil.ReadDir(l.BaseDir)

	if err != nil {
		fmt.Println(err)
	}

	var files []string
	var dirs []string
	for _, file := range dirFs {
		if file.IsDir() {
			dirs = append(dirs, file.Name())
			continue
		}

		files = append(files, file.Name())
	}

	fileCh := make(chan []AnalyzedFile, len(files))
	wg := sync.WaitGroup{}

	wg.Add(len(files))
	for _, fileName := range files {
		go runAnalyze(l, fileName, text, fileCh, &wg)
	}

	wg.Wait()
	close(fileCh)

	for files := range fileCh {
		for _, file := range files {
			file.GetInfo()
		}
	}

	for _, dir := range dirs {
		wg.Add(1)

		loc := NewLocator(l.BaseDir + "/" + dir)

		go runDig(loc, text, &wg)
	}

	wg.Wait()
}

func runAnalyze(locator *Locator, fileName, text string, ch chan<- []AnalyzedFile, wg *sync.WaitGroup) {
	files := locator.Analyze(fileName, text)
	defer wg.Done()

	ch <- files
}

func runDig(locator *Locator, text string, wg *sync.WaitGroup) {
	defer wg.Done()

	locator.Dig(text)
}

// Use Analyze to open and scan the given file, fileName, and assert whether it contains the given string, text. If it does it is accumulated to a slice of AnalyzedFile and later returned
func (l *Locator) Analyze(fileName string, text string) []AnalyzedFile {
	file, err := os.Open(l.BaseDir + "/" + fileName)

	if err != nil {
		fmt.Printf("Path \"%s\" is not valid\n", fileName)
		return []AnalyzedFile{}
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
	var files []AnalyzedFile

	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		line := scanner.Text()

		exists := strings.Contains(line, text)

		if exists {
			resString := strings.Replace(line, text, color.RedString(text), -1)
			files = append(files, AnalyzedFile{FileName: fileName, Content: file, Read: resString, LineNumber: lineNo})
		}

		lineNo++
	}

	return files
}
