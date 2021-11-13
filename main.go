package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

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
		{Type: "README", Filename: filename, Line: 6, Value: ` /* @README
** Headeeeeeer
* Ich bin Grün`},
		{Type: "README", Filename: filename, Line: 7, Value: ` * @README
Irgend n Header
 - irgend n Blödsin
  - Blöd`},
		{Type: "README", Filename: filename, Line: 7, Value: ` * @README
 LOOOOL
 * gdgds
 * gdsfg
  * dfsg
  * gsdg
This is another line`},
	}, nil
}

func ParseCmd() (fileExtensions []string, outputFile string, inputPath string) {
	extVar := flag.String("ext", ".go,.js,.ts,.jsx,.tsx", "")
	outputFileVar := flag.String("out", "", "")
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

	var finder TagFinder = FakeFinder{}
	var loader Loader = FileLoader{fileExtensions}
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
