package core

import (
	"io"
	"path/filepath"

	"github.com/Tiffinger-Thiel-GmbH/AtWhy/finder"
	"github.com/Tiffinger-Thiel-GmbH/AtWhy/loader"
	"github.com/Tiffinger-Thiel-GmbH/AtWhy/processor"
	"github.com/Tiffinger-Thiel-GmbH/AtWhy/tag"
	"github.com/spf13/afero"
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

// AtWhy combines all parts of the application.
// @WHY LINK atwhy_struct_link
// @WHY CODE atwhy_struct_code
type AtWhy struct {
	Loader    Loader
	Finder    loader.TagFinder
	Processor TagProcessor
	Generator Generator
	Writer    io.Writer
}

// @WHY CODE_END

func New(writer io.Writer, gen Generator, projectPath string, extensions []string) (AtWhy, error) {
	abs, err := filepath.Abs(projectPath)
	if err != nil {
		return AtWhy{}, err
	}

	filesystem := afero.NewBasePathFs(afero.NewOsFs(), abs)

	atWhy := AtWhy{
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
	return atWhy, nil
}

func (c AtWhy) Run() error {
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
