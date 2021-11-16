package main

import (
	"io/fs"
	"path/filepath"
	"strings"

	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type FileLoader struct {
	FS             fs.FS
	FileExtensions []string
}

func (fl FileLoader) Load(dir string, finder TagFinder) ([]tag.Raw, error) {
	filesystem := fl.FS
	allTags := make([]tag.Raw, 0)

	// Convert to slash path for compatibility with windows.
	dir = filepath.ToSlash(dir)
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, err
	}

	err = fs.WalkDir(filesystem, absDir[1:], func(path string, info fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		hasFoundExtension := false

		for _, e := range fl.FileExtensions {
			fileName := info.Name()
			if strings.HasSuffix(fileName, e) {
				hasFoundExtension = true
			}
		}

		if !hasFoundExtension {
			return nil
		}

		file, err := filesystem.Open(path)
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
