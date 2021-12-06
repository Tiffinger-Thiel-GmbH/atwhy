package core

import (
	"io"
	"path/filepath"

	"github.com/spf13/afero"
	"gitlab.com/tiffinger-thiel/crazydoc/finder"
	"gitlab.com/tiffinger-thiel/crazydoc/loader"
	"gitlab.com/tiffinger-thiel/crazydoc/processor"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type Loader interface {
	Load(finder loader.TagFinder) (allTags []tag.Raw, err error)
}

type TagProcessor interface {
	Process(tags []tag.Raw) ([]tag.Tag, error)
}

type Generator interface {
	Generate(tags []tag.Tag, writer io.Writer) error
}

// CrazyDoc combines all parts of the application.
// @DOC LINK crazydoc_struct_link
// @DOC CODE crazydoc_struct_code
type CrazyDoc struct {
	Loader    Loader
	Finder    loader.TagFinder
	Processor TagProcessor
	Generator Generator
	Writer    io.Writer
}

// @DOC CODE_END

func New(writer io.Writer, gen Generator, projectPath string, extensions []string) (CrazyDoc, error) {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return CrazyDoc{}, err
	}

	filesystem := afero.NewBasePathFs(afero.NewOsFs(), abs)

	crazyDoc := CrazyDoc{
		Finder: &finder.Finder{
			BlockCommentStarts: []string{"/*"},
			BlockCommentEnds:   []string{"*/"},
			LineCommentStarts:  []string{"//"},
		},
		Loader: loader.File{
			FS:             filesystem,
			FileExtensions: extensions,
		},
		Processor: processor.Processor{
			TagFactories: []tag.Factory{
				tag.Doc,
				tag.Code,
				tag.Link,
			},
		},
		Generator: gen,
		Writer:    writer,
	}
	return crazyDoc, nil
}

func (c CrazyDoc) Run() error {
	tags, err := c.Loader.Load(c.Finder)
	if err != nil {
		return err
	}

	processed, err := c.Processor.Process(tags)
	if err != nil {
		return err
	}

	return c.Generator.Generate(processed, c.Writer)
}
