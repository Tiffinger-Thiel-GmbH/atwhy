package generator

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/spf13/afero"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type DocTemplate struct {
	ID    string
	Name  string
	Value string

	template *template.Template
}

// LoadDocTemplates from the given folder.
// If allowList is not nil, it only loads the filenames included in the allowList.
//
//  @DOC doc_template
//  The templates should be normal markdown files.
//  The first line has to be the name of the template (used for example for the navigation in the html-generator).
//
//  You can access a tag called `\@DOC example_tag` using
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
			file, err := sysfs.Open(path)
			if err != nil {
				return err
			}
			defer file.Close()

			tplData, err := ioutil.ReadAll(file)
			if err != nil {
				return err
			}

			// Windows compatibility:
			tplData = bytes.ReplaceAll(tplData, []byte("\r\n"), []byte("\n"))
			splitted := bytes.SplitN(tplData, []byte("\n"), 2)
			var name string
			var body string
			if len(splitted) < 1 {
				return errors.New("the template must have at least one line containing the name")
			}

			name = string(splitted[0])

			if len(splitted) >= 2 {
				body = string(splitted[1])
			}

			id := md5.Sum([]byte(filepath.ToSlash(path)))
			tpl, err := template.New(name).Parse(body)

			if err != nil {
				return err
			}

			res = append(res, DocTemplate{
				ID:       "page-" + hex.EncodeToString(id[:]),
				Name:     name,
				Value:    body,
				template: tpl,
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Sort by name.
	sort.Slice(res, func(i, j int) bool {
		return res[i].Name < res[j].Name
	})

	return res, nil
}

// Execute the template
func (t DocTemplate) Execute(tagMap map[string]tag.Tag, writer io.Writer) error {
	data := struct {
		Tag map[string]tag.Tag
		Now string
	}{
		Tag: tagMap,
		Now: time.Now().Format(time.RFC822Z),
	}

	return t.template.Execute(writer, data)
}
