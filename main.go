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
	tags, err := loader.Load(path, finder)
	if err != nil {
		panic(err)
	}
	fmt.Println(tags)
}
