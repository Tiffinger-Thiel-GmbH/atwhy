package generator

import (
	"crypto/md5"
	"encoding/hex"
	"io"
	"io/fs"
	"io/ioutil"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/spf13/afero"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type DocTemplate struct {
	ID    string
	Name  string
	Value string
}

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

			id := md5.Sum([]byte(filepath.ToSlash(path)))
			name, err := filepath.Rel(abs, path)
			if err != nil {
				return err
			}
			res = append(res, DocTemplate{
				ID: "page-" + hex.EncodeToString(id[:]),

				// TODO: could be read from first line in the template...
				Name:  strings.TrimSuffix(name, ".tpl.md"),
				Value: string(tplData),
			})
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return res, nil
}

func (t DocTemplate) Execute(tagMap map[string]tag.Tag, writer io.Writer) error {
	tpl, err := template.New(t.Name).Parse(t.Value)
	if err != nil {
		return err
	}

	data := struct {
		Tag map[string]tag.Tag
	}{
		Tag: tagMap,
	}

	return tpl.Execute(writer, data)
}
