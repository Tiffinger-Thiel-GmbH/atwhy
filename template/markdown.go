package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"gopkg.in/yaml.v2"

	"github.com/spf13/afero"
)

const templateSuffix = ".tpl.md"

var (
	ErrMissingBody = errors.New("if the first line is '---' you have to include a yaml header as described in the atwhy readme")
)

// @WHY doc_template_header_1
// Each template may have a yaml Header.
// Example with all possible fields:
// ```markdown
// ---
// # Some metadata which may be used for the generation.
// meta:
//   # The title is used for the served html to e.g. generate a menu and add page titles.
//   title: Readme # default: the template filename
//
// # Additional configuration for the `atwhy serve` command.
// server:
//   index: true # default: false
// ---
// # Your Markdown starts here
//
// ## Foo
// bar
// ```
// (Note: VSCode supports the header automatically.)

type Header struct {
	// Meta contains additional data which can be used by the generators.
	// It is also available inside the template for example with
	//  {{ .Meta.Title }}
	Meta MetaData `yaml:"meta"`

	Server ServerData `yaml:"server"`
}

type MetaData struct {
	// Title is for example used in the html generator to create the navigation buttons.
	// If not set, it will default to the template file-name (excluding .tpl.md)
	Title string `yaml:"title"`
}

type ServerData struct {
	// Index defines if this template should be used as "index.html".
	// Note that there can only be one page in each folder which is the index.
	Index bool `yaml:"index"`
}

// Markdown
//
// @WHY doc_template
// The templates should be markdown files with a yaml header for metadata.
//
// You can access a tag called `\@WHY example_tag` using
//
//	# Example
//	{{ .Escape
//	"{{ .Tag.example_tag }}"
//	}}
//
// Note: This uses the Go templating engine.
// Therefor you can use the [Go templating syntax](https://learn.hashicorp.com/tutorials/nomad/go-template-syntax?in=nomad/templates).
type Markdown struct {
	ID                string
	ProjectPathPrefix string
	Name              string
	Path              string
	Header            Header

	template *template.Template
	tagMap   map[string]tag.Tag
}

func readTemplate(sysfs afero.Fs, projectPathPrefix string, path string, tags mappedTags) (Markdown, error) {
	file, err := sysfs.Open(path)
	if err != nil {
		return Markdown{}, err
	}
	defer file.Close()

	tplData, err := io.ReadAll(file)
	if err != nil {
		return Markdown{}, err
	}

	// Windows compatibility:
	tplData = bytes.ReplaceAll(tplData, []byte("\r\n"), []byte("\n"))

	var body string
	header := Header{}

	// No Header exists because the first line was no "---"
	if bytes.HasPrefix(tplData, []byte("---\n")) {
		// Extract the Header:
		splitted := bytes.SplitN(tplData, []byte("---\n"), 3)

		if len(splitted) == 3 {
			body = string(splitted[2])

			err = yaml.Unmarshal(splitted[1], &header)
			if err != nil {
				return Markdown{}, err
			}
		} else {
			return Markdown{}, ErrMissingBody
		}
	} else {
		body = string(tplData)
	}

	filename := filepath.Base(path)

	id := md5.Sum([]byte(filepath.ToSlash(path)))
	tpl, err := template.New(filename).Parse(body)

	if err != nil {
		return Markdown{}, err
	}

	if header.Meta.Title == "" {
		header.Meta.Title = strings.TrimSuffix(filename, templateSuffix)
	}

	markdownTemplate := Markdown{
		ID:                "page-" + hex.EncodeToString(id[:]),
		ProjectPathPrefix: projectPathPrefix,
		Name:              strings.TrimSuffix(filepath.Base(path), templateSuffix),
		Path:              filepath.Dir(path),

		Header:   header,
		template: tpl,
		tagMap:   tags,
	}

	return markdownTemplate, nil
}

type data struct {
	Tag           map[string]tag.Tag
	Meta          MetaData
	Now           string
	projectPrefix string

	isPostprocessing bool
}

func (d data) Project(file string) string {
	return filepath.Join(d.projectPrefix, file)
}

// Escape provides escaping of {{ and }} which works in the
// first execution and also in the post-process-execution.
//
// The default Go-Template way {{"{{   }}"}} works only in the
// tags and not in the templates due to the preprocessing.
//
// @WHY doc_template_escape_tag
// __What if `{{"{{"}}` or `{{"}}"}}` is needed in the documentation?__
// You can wrap them with `{{ .Escape "{{.Escape \"...\"}}" }}`.
// E.g.: `{{ .Escape "{{ .Escape \"\\\"{{\\\"  and  \\\"}}\\\"\" }}" }}`
// Results in this markdown text: `{{ .Escape "\"{{\" and \"}}\"" }}`
//
// __Note:__ You need to escape `"` with `\"`.
//
// (The official Go-Template way `{{ .Escape "{{ \"{{ -- }}\" }}" }}` doesn't work in all cases with atwhy. `.Escape` works always.)
func (d data) Escape(value string) string {
	if !d.isPostprocessing {
		value = strings.ReplaceAll(value, `\`, `\\`)
		value = strings.ReplaceAll(value, `"`, `\"`)
		value = `{{"` + value + `"}}`
	}
	return value
}

// Execute the template
func (t Markdown) Execute(writer io.Writer) error {

	// @WHY doc_template_possible_tags
	// __Possible template values are:__
	// * Any Tag from the project: `{{"{{ .Tag.example_tag }}"}}`
	// * Current Datetime: `{{"{{ .Now }}"}}`
	// * Metadata from the yaml header: `{{"{{ .Meta.Title }}"}}`
	// * Conversion of links to project-files (also in serve-mode): `{{"{{ .Project \"my/file/in/the/project.go\" }}"}}`
	//   You need to use that if you want to generate links to actual files in your project.
	//   This can also be used for pictures: `{{ .Escape "![aPicture]({{ .Project \"path/to/the/picture.jpg\" }})" }}`

	d := data{
		Tag:  t.tagMap,
		Now:  time.Now().Format(time.RFC822Z), // TODO: add this as function instead of as value to be able to pass any format.
		Meta: t.Header.Meta,

		projectPrefix: t.ProjectPathPrefix,
	}

	buf := bytes.NewBufferString("")

	// First just execute the template.
	err := t.template.Execute(buf, d)
	if err != nil {
		return err
	}

	// And then execute the postprocessing template.
	// E.g. it can process the {{ .Project }} even if the links are inside the tags.
	postProcessTemplate, err := template.New("postProcessing.md").Parse(buf.String())
	if err != nil {
		return err
	}
	// Do not allow tags in this step as it would create bad edge cases.
	d.Tag = map[string]tag.Tag{}
	d.isPostprocessing = true
	err = postProcessTemplate.Execute(writer, d)

	return err
}
