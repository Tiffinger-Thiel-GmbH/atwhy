package main

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

type TagType string

var (
	TagWhy      TagType = "WHY"
	TagReadme   TagType = "README"
	TagFileLine TagType = "FILELINE"
	TagFlag     TagType = "FLAG"
)

type Tag struct {
	Type     TagType
	Filename string
	Line     int
	Value    string
}

type FileLoader interface {
	Load(dir string, finder TagFinder) (allTags []Tag, err error)
}

type TagFinder interface {
	Find(filename string, reader io.Reader) (tags []Tag, err error)
	// SaveByTag()
	// scan()
	// findTag()
	// saveByTag()
}

type ProcessedTag struct {
	Type     TagType
	Value    string
	Children []ProcessedTag
}

type TagProcessor interface {
	Process(tags []Tag) ([]ProcessedTag, error)
}

type Generator interface {
	Generate(tags []ProcessedTag, writer io.Writer) error
}

type FakeFinder struct {
}

func (ff FakeFinder) Find(filename string, reader io.Reader) (tags []Tag, err error) {
	return []Tag{
		{Type: TagReadme, Filename: filename, Line: 5, Value: "jdfglh"},
	}, nil
}

func main() {
	ext := flag.String("ext", "", "")
	flag.Parse()
	path := flag.Arg(0)
	fileExtensions := strings.Split(*ext, ",")

	var finder TagFinder = FakeFinder{}
	var loader FileLoader = Loader{fileExtensions}
	_, err := loader.Load(path, finder)
	if err != nil {
		panic(err)
	}

	var mockedTags []ProcessedTag = []ProcessedTag{{Type: TagReadme, Value: "# Config \n ## API \n Test 123", Children: []ProcessedTag{{Type: TagFileLine, Value: "[13:1]", Children: nil}}}}

	var g Generate = Generate{}
	b := strings.Builder{}
	fmt.Println(g.Generate(mockedTags, &b))
	fmt.Println(b.String())
}
