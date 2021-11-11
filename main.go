package main

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/spf13/afero"
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
	Find(filename string, reader io.Reader) error
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
	Generate(tags []ProcessedTag) (io.Reader, error)
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
