package main

import (
	"fmt"
	"io"
)

type FileLoader interface {
}

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

type TagFinder interface {
	Find(filename string, reader io.Reader) ([]Tag, error)
	// SaveByTag()
	// scan()
	// findTag()
	// saveByTag()
}

type ProcessedTag struct {
	Type  TagType
	Value string
}

type TagProcessor interface {
	Process(tags []Tag) ([]ProcessedTag, error)
}

type Generator interface {
	Generate(tags []ProcessedTag) (io.Reader, error)
}

func main() {
	fmt.Println("Hello World")
}
