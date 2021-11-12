package main

import (
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
	"io/fs"
	"strings"

	"github.com/spf13/afero"
)

type Loader struct {
	FileExtensions []string
}

func (l Loader) Load(dir string, finder TagFinder) ([]tag.Raw, error) {
	var AppFs = afero.NewOsFs()
	allTags := make([]tag.Raw, 0)

	err := afero.Walk(AppFs, dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		hasFoundExtension := false

		for _, e := range l.FileExtensions {
			fileName := info.Name()
			if strings.HasSuffix(fileName, e) {
				hasFoundExtension = true
			}
		}

		if !hasFoundExtension {
			return nil
		}

		file, err := AppFs.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		tags, err := finder.Find(path, file)
		if err != nil {
			return err
		}
		allTags = append(allTags, tags...)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return allTags, nil
}
