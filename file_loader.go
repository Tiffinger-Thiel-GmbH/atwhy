package main

import (
	"io/fs"
	"strings"

	"github.com/spf13/afero"
)

type Loader struct {
	FileExtensions []string
}

func (l Loader) Load(dir string, finder TagFinder) ([]Tag, error) {
	var AppFs = afero.NewOsFs()
	allTags := []Tag{}

	err := afero.Walk(AppFs, dir, func(path string, info fs.FileInfo, err error) error {

		if info.IsDir() {
			return nil
		}

		hasFoundExtension := true

		for _, e := range l.FileExtensions {
			fileName := info.Name()
			if !strings.HasSuffix(fileName, e) {
				hasFoundExtension = false
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
