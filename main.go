package main

import (
	"flag"
	"fmt"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
	"io"
	"os"
	"strings"
)

// @README 10
// What is CrazyDoc
// CrazyDoc can be used to generate a documentation out of comments in the code.
// That way you can for example describe all available options in the same file
// where they are coded. A developer therefore doesn't have to know exactly where
// the information has to be documented because it is just in the same file.
//
// The same applies to architectural decisions, which can be documented, where its
// actually done.
// --> __Single source of truth__ also for documentation!

// @README 20
// Distribute
// # Prerequisites
// * Go 1.17
//
// # Build
// Run `go build .`
//

type Loader interface {
	Load(dir string, finder TagFinder) (allTags []tag.Raw, err error)
}

type TagFinder interface {
	Find(filename string, reader io.Reader) (tags []tag.Raw, err error)
	// SaveByTag()
	// scan()
	// findTag()
	// saveByTag()
}

type TagProcessor interface {
	Process(tags []tag.Raw) ([]tag.Tag, error)
}

type Generator interface {
	Generate(tags []tag.Tag, writer io.Writer) error
}

type FakeFinder struct {
}

func (ff FakeFinder) Find(filename string, reader io.Reader) (tags []tag.Raw, err error) {
	return []tag.Raw{
		{Type: "FILELINK", Filename: filename, Line: 5, Value: `// @FILELINK`},
		{Type: "README", Filename: filename, Line: 6, Value: ` /* @README 10
** Headeeeeeer
* Ich bin Grün`},
		{Type: "README", Filename: filename, Line: 7, Value: ` * @README 5
Irgend n Header
 - irgend n Blödsin
  - Blöd`},
		{Type: "README", Filename: filename, Line: 7, Value: ` * @README 20
 LOOOOL
 * gdgds
 * gdsfg
  * dfsg
  * gsdg
This is another line`},
	}, nil
}

func ParseCmd() (fileExtensions []string, outputFile string, inputPath string) {
	// @README 20
	// Usage
	// Just run `crazydoc [OPTIONS]... [PROJECT_ROOT]`.
	// To get all possible file extensions just run `crazydoc -help`

	flag.Usage = func() {
		_, _ = fmt.Fprintf(flag.CommandLine.Output(), "Usage: %s [OPTIONS]... [PROJECT_ROOT]\n", os.Args[0])
		flag.PrintDefaults()
	}
	extVar := flag.String("ext", ".go,.js,.ts,.jsx,.tsx", "comma separated list of file extensions to search for")
	outputFileVar := flag.String("out", "", "ouptut file \nshould be a .md file")
	flag.Parse()

	inputPath = flag.Arg(0)
	fileExtensions = strings.Split(*extVar, ",")
	if inputPath == "" {
		inputPath = "."
	}

	return fileExtensions, *outputFileVar, inputPath
}

func main() {
	
	fileExtensions, outputFile, inputPath := ParseCmd()

	var finder TagFinder = Finder{}
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	var loader Loader = FileLoader{
		FS:             os.DirFS(currentDir),
		FileExtensions: fileExtensions,
	}
	var processor TagProcessor = Processor{
		cleaners: []Cleaner{
			SlashStarCleaner{},
		},
		tagFactories: []tag.Factory{
			tag.Why,
			tag.Readme,
			tag.FileLink,
		},
	}
	var generator Generator = MarkdownGenerator{}

	writer := os.Stdout
	if outputFile != "" {
		file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		writer = file
	}

	crazyDoc := CrazyDoc{
		Finder:    finder,
		Loader:    loader,
		Processor: processor,
		Generator: generator,
		Writer:    writer,
	}

	if err := crazyDoc.Run(inputPath); err != nil {
		panic(err)
	}
}
