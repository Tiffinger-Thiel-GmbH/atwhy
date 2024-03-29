package core

import (
	"html/template"
	"io"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/finder"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/loader"
	mdTemplate "github.com/Tiffinger-Thiel-GmbH/atwhy/template"
	"github.com/spf13/afero"
)

type Loader interface {
	Load(finder loader.TagFinder) (allTags []tag.Raw, err error)
}

type TagProcessor interface {
	Process(tags []tag.Raw) ([]tag.Tag, error)
}

type Generator interface {
	Generate(markdownTemplate mdTemplate.Markdown, writer io.Writer) error

	// Ext returns the file extension which should be used for the generated files.
	Ext() string
}

type TemplateLoader interface {
	Load(tags []tag.Tag) ([]mdTemplate.Markdown, error)
}

// @WHY atwhy_interfaces
// * `Loader` loads files from a given path.
// * `loader.TagFinder` reads the file and returns all lines which are part of a found tag. It Does not process the raw lines.
// * `TagFactories` convert the raw tags from the `TagFinder` and generates final Tags out of them.
// * `TemplateLoader` loads the templates from the `template` folder to pass them the generator.
// * `Generator` is responsible for postprocessing the tags and output the final file. which it just writes to the
// passed `Writer`.
//
// So the workflow is:
// Loader -> TagFinder = tagList []tag.Raw tagList -> TagProcessor -> TemplateLoader -> Generator -> Writer

// AtWhy combines all parts of the application.
// @WHY LINK atwhy_struct_link
// @WHY CODE atwhy_struct_code
type AtWhy struct {
	Loader         Loader
	Finder         loader.TagFinder
	TagFactories   []tag.Factory
	Generator      Generator
	TemplateLoader TemplateLoader

	projectPath       string
	projectPathPrefix string
	pageTemplate      *template.Template
}

// @WHY CODE_END

func New(gen Generator, projectPath string, projectPathPrefix string, templateFolder string, extensions []string, commentConfig map[string]finder.CommentConfig) (AtWhy, error) {
	filesystem := afero.NewBasePathFs(afero.NewOsFs(), projectPath)
	templateFS := afero.NewBasePathFs(filesystem, templateFolder)

	atwhy := AtWhy{
		Finder: &finder.Finder{
			CommentConfig: commentConfig,
		},
		Loader: loader.File{
			FS:             filesystem,
			FileExtensions: extensions,
		},
		TagFactories: []tag.Factory{
			tag.Doc,
			tag.Code,
			tag.ProjectLink,
		},
		Generator: gen,
		TemplateLoader: mdTemplate.Loader{
			FS:                templateFS,
			ProjectPathPrefix: projectPathPrefix,
		},

		projectPath:       projectPath,
		projectPathPrefix: projectPathPrefix,
	}

	err := atwhy.initPageTemplate()
	if err != nil {
		return AtWhy{}, err
	}

	return atwhy, nil
}

func (a *AtWhy) Load() ([]mdTemplate.Markdown, error) {
	tags, err := a.Loader.Load(a.Finder)
	if err != nil {
		return nil, err
	}

	var processed []tag.Tag

	for i := range tags {
		t := tags[i]

		for _, factory := range a.TagFactories {
			newTag, err := factory(t)
			if err != nil {
				return nil, err
			}
			if newTag == nil {
				continue
			}

			processed = append(processed, newTag)
		}
	}

	return a.TemplateLoader.Load(processed)
}

func (a *AtWhy) Generate(template mdTemplate.Markdown, writer io.Writer) error {
	return a.Generator.Generate(template, writer)
}
