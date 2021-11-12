package main

import (
	"flag"
	"fmt"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
	"io"
	"strings"
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
}

type FakeFinder struct {
}

func (ff FakeFinder) Find(filename string, reader io.Reader) (tags []tag.Raw, err error) {
	return []tag.Raw{
		{Type: "FILELINK", Filename: filename, Line: 5, Value: `// @FILELINK`},
		{Type: "README", Filename: filename, Line: 6, Value: ` /* @README`},
		{Type: "README", Filename: filename, Line: 7, Value: ` * @README`},
		{Type: "README", Filename: filename, Line: 7, Value: ` * @README
LOOOOL
 * gdgds
 * gdsfg
  * dfsg
  * gsdg
This is another line`},
		{Type: "README", Filename: filename, Line: 8, Value: "jdfglh"},
	}, nil
}

func main() {
	ext := flag.String("ext", ".go,.js,.ts,.jsx,.tsx", "")
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
	fmt.Println(tags)

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
	fmt.Println(processed)
}
