package main

import (
	"fmt"
	"io"
)

type FileLoader interface {
}

type TagType string

var (
	Why TagType = "WHY"
	Readme TagType = "README"
	FileLine TagType = "FILELINE"
)

type Tag struct {
	Filename string
	Line int
	Value string
}

type TagFinder interface {
	Find(filename string, reader io.Reader) ([]Tag, error)
	// SaveByTag()
	// scan()
	// findTag()
	// saveByTag()
}

type CommentCleaner interface {
}

type Generator interface {
}

type Writer interface {
}

func main() {
	fmt.Println("Hello World")

}
