package generator

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/tag"
	"github.com/spf13/afero"
	"gopkg.in/yaml.v2"
)

const templateSuffix = ".tpl.md"

// @WHY doc_template_header_1
// Each template can have a yaml header with the following fields:

// @WHY CODE doc_template_header_struct
type TemplateHeader struct {
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

type DocTemplate struct {
	ID     string
	Value  string
	Header TemplateHeader

	template *template.Template
}

func readTemplate(sysfs afero.Fs, path string) (DocTemplate, error) {
	file, err := sysfs.Open(path)
	if err != nil {
		return DocTemplate{}, err
	}
	defer file.Close()

	tplData, err := ioutil.ReadAll(file)
	if err != nil {
		return DocTemplate{}, err
	}

	// Windows compatibility:
	tplData = bytes.ReplaceAll(tplData, []byte("\r\n"), []byte("\n"))

	// Extract the header:
	splitted := bytes.SplitN(tplData, []byte("---\n"), 2)
	header := TemplateHeader{}
	var body string

	// No header exists:
	if len(splitted) == 1 {
		body = string(splitted[0])
	}

	if len(splitted) >= 2 {
		body = string(splitted[1])
		err = yaml.Unmarshal(splitted[0], &header)
		if err != nil {
			return DocTemplate{}, err
		}
	}

	filename := filepath.Base(path)

	id := md5.Sum([]byte(filepath.ToSlash(path)))
	tpl, err := template.New(filename).Parse(body)

	if err != nil {
		return DocTemplate{}, err
	}

	if header.Meta.Title == "" {
		header.Meta.Title = strings.TrimSuffix(filename, templateSuffix)
	}

	return DocTemplate{
		ID:       "page-" + hex.EncodeToString(id[:]),
		Value:    body,
		Header:   header,
		template: tpl,
	}, nil
}

// LoadDocTemplates from the given folder.
// If allowList is not nil, it only loads the filenames included in the allowList.
//
//  @WHY doc_template
//  The templates should be normal markdown files.
//  The first line has to be the name of the template (used for example for the navigation in the html-generator).
//
//  You can access a tag called `\@WHY example_tag` using
//  ```text
//  # Example
//  {{ .Tag.example_tag }}
//  ```
//
//  Note: This is basically the syntax of the Go templating engine.
//  Therefor you can use the [Go templating syntax](https://learn.hashicorp.com/tutorials/nomad/go-template-syntax?in=nomad/templates).
func LoadDocTemplates(sysfs afero.Fs, folder string, allowList []string) ([]DocTemplate, error) {
	var res []DocTemplate
	abs, err := filepath.Abs(folder)
	if err != nil {
		return nil, err
	}
	err = afero.Walk(sysfs, abs, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		// Filter for names in the allow list if it is not nil.
		if len(allowList) > 0 {
			var allowed bool
			for _, name := range allowList {
				if info.Name() == name+".tpl.md" || info.Name() == name {
					allowed = true
					break
				}
			}

			if !allowed {
				return nil
			}
		}

		if strings.HasSuffix(path, ".tpl.md") {
			newTpl, err := readTemplate(sysfs, path)
			if err != nil {
				return err
			}

			res = append(res, newTpl)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by the title.
	sort.Slice(res, func(i, j int) bool {
		return res[i].Header.Meta.Title < res[j].Header.Meta.Title
	})

	return res, nil
}

// Execute the template
func (t DocTemplate) Execute(tagMap map[string]tag.Tag, writer io.Writer) error {
	data := struct {
		Tag  map[string]tag.Tag
		Meta MetaData
		Now  string
	}{
		Tag: tagMap,
		Now: time.Now().Format(time.RFC822Z),
	}

	return t.template.Execute(writer, data)
}
