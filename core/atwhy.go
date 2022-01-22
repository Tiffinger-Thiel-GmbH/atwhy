package core

import (
	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/finder"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/loader"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/processor"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/template"
	"github.com/spf13/afero"
	"io"
)

type Loader interface {
	Load(finder loader.TagFinder) (allTags []tag.Raw, err error)
}

type TagProcessor interface {
	Process(tags []tag.Raw) ([]tag.Tag, error)
}

type Generator interface {
	Generate(markdownTemplate template.MarkdownTemplate, writer io.Writer) error

	// Ext returns the file extension which should be used for the generated files.
	Ext() string
}

type TemplateLoader interface {
	Load(tags []tag.Tag) ([]template.MarkdownTemplate, error)
}

// AtWhy combines all parts of the application.
// @WHY LINK atwhy_struct_link
// @WHY CODE atwhy_struct_code
type AtWhy struct {
	Loader         Loader
	Finder         loader.TagFinder
	Processor      TagProcessor
	Generator      Generator
	TemplateLoader TemplateLoader
}

// @WHY CODE_END

func New(gen Generator, projectPath string, templateFolder string, extensions []string) (AtWhy, error) {
	filesystem := afero.NewBasePathFs(afero.NewOsFs(), projectPath)
	templateFS := afero.NewBasePathFs(filesystem, templateFolder)

	atwhy := AtWhy{
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
		TemplateLoader: template.Loader{
			FS: templateFS,
		},
	}
	return atwhy, nil
}

func (c *AtWhy) Load() ([]template.MarkdownTemplate, error) {
	tags, err := c.Loader.Load(c.Finder)
	if err != nil {
		return nil, err
	}

	processedTags, err := c.Processor.Process(tags)
	if err != nil {
		return nil, err
	}

	return c.TemplateLoader.Load(processedTags)
}

func (c *AtWhy) Generate(template template.MarkdownTemplate, writer io.Writer) error {
	return c.Generator.Generate(template, writer)
}
