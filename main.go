package main

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/spf13/afero"
)

type TagType string

var (
	Why      TagType = "WHY"
	Readme   TagType = "README"
	FileLine TagType = "FILELINE"
)

type Tag struct {
	Filename string
	Line     int
	Value    string
}

type FileLoader interface {
	Load(dir string, finder TagFinder) (allTags []Tag, err error)
}

type TagFinder interface {
	Find(filename string, reader io.Reader) error
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
	var AppFs = afero.NewOsFs()

	err := afero.Walk(AppFs, "../GoKt", func(path string, info fs.FileInfo, err error) error {
		fmt.Println(path)
		return nil
	})
	if err != nil {
		panic(err)
	}
}
