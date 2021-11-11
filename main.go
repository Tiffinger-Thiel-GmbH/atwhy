package main

import (
	"fmt"
	"io"
)

type FileLoader interface {
}

type TagType string

var (
	Why      TagType = "WHY"
	Readme   TagType = "README"
	FileLine TagType = "FILELINE"
	Flag     TagType = "FLAG"
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
	Value interface{}
}

type TagProcessor interface {
	Process(tags []Tag) ([]ProcessedTag, error)
}

type Generator interface {
}

func main() {
	fmt.Println("Hello World")
}
