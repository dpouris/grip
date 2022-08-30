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
	Ok         bool
}

func (f *AnalyzedFile) GetInfo() {
	// return file info
	fileName := color.YellowString(f.FileName)
	lineNo := color.GreenString(fmt.Sprintf("%v", f.LineNumber))
	fmt.Printf("\nFilename: %s\nLine Number: %s\nFound: %v\n", fileName, lineNo, f.Read)
}

// We'll make a function NewLocator that takes the directory and will init a locator
func NewLocator(dir string) *Locator {
	loc := Locator{BaseDir: dir}
	return &loc
}

// We'll make a function that will search through the directory given a desired string
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

	fileCh := make(chan *AnalyzedFile, len(files))
	wg := sync.WaitGroup{}

	wg.Add(len(files))
	for _, fileName := range files {
		go runAnalyze(l, fileName, text, fileCh, &wg)
	}

	wg.Wait()
	close(fileCh)

	for file := range fileCh {
		file.GetInfo()
	}

	for _, dir := range dirs {
		wg.Add(1)

		loc := NewLocator(l.BaseDir + "/" + dir)

		go runDig(loc, text, &wg)
	}

	wg.Wait()
}

func runAnalyze(locator *Locator, fileName, text string, ch chan<- *AnalyzedFile, wg *sync.WaitGroup) {
	file := locator.Analyze(fileName, text)
	defer wg.Done()

	if !file.Ok {
		// Handle not found
		// fmt.Printf("Not found in %s\n", fileName)
		return
	}

	ch <- file
}

func runDig(locator *Locator, text string, wg *sync.WaitGroup) {
	defer wg.Done()

	locator.Dig(text)
}

// We'll make a function called Analyze that opens a file and asserts whether a given string is inside that file and returns the opened file ptr along with a bool value indicating the validity of the file
func (l *Locator) Analyze(fileName string, text string) *AnalyzedFile {
	file, err := os.Open(l.BaseDir + "/" + fileName)

	if err != nil {
		fmt.Printf("Path \"%s\" is not valid\n", fileName)
		return &AnalyzedFile{}
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
	lineNo := 0

	for scanner.Scan() {

		if err := scanner.Err(); err != nil {
			panic(err)
		}

		line := scanner.Text()

		exists := strings.Contains(line, text)

		if exists {
			resString := strings.Replace(line, text, color.RedString(text), -1)
			return &AnalyzedFile{FileName: fileName, Content: file, Read: resString, LineNumber: lineNo, Ok: true}
		}

		lineNo++
	}

	return &AnalyzedFile{}
}
