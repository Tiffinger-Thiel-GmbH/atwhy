package main

import (
	"flag"
	"io"
	"os"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type FileLoader interface {
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

func main() {
	ext := flag.String("ext", ".go,.js,.ts,.jsx,.tsx", "")
	outputFile := flag.String("out", "", "")
	flag.Parse()
	path := flag.Arg(0)
	fileExtensions := strings.Split(*ext, ",")
	if path == "" {
		path = "."
	}

	var finder TagFinder = FakeFinder{}
	var loader FileLoader = Loader{fileExtensions}
	tags, err := loader.Load(path, finder)
	if err != nil {
		panic(err)
	}

	processor := Processor{
		cleaners: []Cleaner{
			SlashStarCleaner{},
		},
		tagFactories: []tag.Factory{
			tag.Why,
			tag.Readme,
			tag.FileLink,
		},
	}

	processed, err := processor.Process(tags)
	if err != nil {
		panic(err)
	}

	var g Generate = Generate{}
	// TODO stdout when no -out
	if *outputFile != "" {
		file, err := os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			panic(err)
		}
		defer file.Close()

		g.Generate(processed, file)
	} else {
		g.Generate(processed, os.Stdout)
	}
}
