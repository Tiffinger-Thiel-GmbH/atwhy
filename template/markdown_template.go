package template

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"io"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const templateSuffix = ".tpl.md"

// @WHY doc_template_header_1
// Each template has a yaml Header with the following fields:

// @WHY CODE doc_template_header_struct

type Header struct {
	// Meta contains additional data which can be used by the generators.
	// It is also available inside the template for example with
	//  {{ .Meta.Title }}
	Meta MetaData `yaml:"meta"`
}

type MetaData struct {
	// Title is for example used in the html generator to create the navigation buttons.
	// If not set, it will default to the template file-name (excluding .tpl.md)
	Title string `yaml:"title"`
}

// @WHY CODE_END

// @WHY doc_template_header_2
// The header is separated from the markdown by using a line with three `-` and a newline.
// Example:
// ```md
// meta:
//  title: Readme
// ---
// # Your Markdown
//
// ## Foo
// bar
// ```
//

// MarkdownTemplate
//
// @WHY doc_template
// The templates should be markdown files with a yaml header for metadata.
//
// You can access a tag called `\@WHY example_tag` using
//  ```text
//  # Example
//  {{ .Tag.example_tag }}
//  ```
//
// Note: This uses the Go templating engine.
// Therefor you can use the [Go templating syntax](https://learn.hashicorp.com/tutorials/nomad/go-template-syntax?in=nomad/templates).
type MarkdownTemplate struct {
	ID     string
	Name   string
	Path   string
	Value  string
	Header Header

	template *template.Template
	tagMap   map[string]tag.Tag
}

func readTemplate(sysfs afero.Fs, path string, tagMap map[string]tag.Tag) (MarkdownTemplate, error) {
	file, err := sysfs.Open(path)
	if err != nil {
		return MarkdownTemplate{}, err
	}
	defer file.Close()

	tplData, err := ioutil.ReadAll(file)
	if err != nil {
		return MarkdownTemplate{}, err
	}

	// Windows compatibility:
	tplData = bytes.ReplaceAll(tplData, []byte("\r\n"), []byte("\n"))

	// Extract the Header:
	splitted := bytes.SplitN(tplData, []byte("---\n"), 2)
	header := Header{}
	var body string

	// No Header exists:
	if len(splitted) == 1 {
		body = string(splitted[0])
	}

	if len(splitted) >= 2 {
		body = string(splitted[1])
		err = yaml.Unmarshal(splitted[0], &header)
		if err != nil {
			return MarkdownTemplate{}, err
		}
	}

	filename := filepath.Base(path)

	id := md5.Sum([]byte(filepath.ToSlash(path)))
	tpl, err := template.New(filename).Parse(body)

	if err != nil {
		return MarkdownTemplate{}, err
	}

	if header.Meta.Title == "" {
		header.Meta.Title = strings.TrimSuffix(filename, templateSuffix)
	}

	return MarkdownTemplate{
		ID:    "page-" + hex.EncodeToString(id[:]),
		Name:  strings.TrimSuffix(filepath.Base(path), templateSuffix),
		Path:  filepath.Dir(path),
		Value: body,

		Header:   header,
		template: tpl,
		tagMap:   tagMap,
	}, nil
}

// Execute the template
func (t MarkdownTemplate) Execute(writer io.Writer) error {

	data := struct {
		Tag  map[string]tag.Tag
		Meta MetaData
		Now  string
	}{
		Tag: t.tagMap,

		// @WHY doc_template_possible_tags
		// Possible template values are:
		// * Any Tag from the project: `{{ .Tag.example_tag }}`
		// * Current Datetime: `{{ .Now }}`

		Now: time.Now().Format(time.RFC822Z),
	}

	return t.template.Execute(writer, data)
}
