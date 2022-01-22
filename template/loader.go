package template

import (
	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/spf13/afero"
	"io/fs"
	"sort"
	"strings"
)

type Loader struct {
	FS afero.Fs
}

func createTagMap(tags []tag.Tag, tpl Markdown) map[string]tag.Tag {
	tagMap := make(map[string]tag.Tag)

	for _, t := range tags {
		tagMap[t.Placeholder()] = t.WithContext(tpl)
	}

	return tagMap
}

// Load templates from the Loader.FS.
func (l Loader) Load(tags []tag.Tag) ([]Markdown, error) {
	var res []Markdown

	err := afero.Walk(l.FS, "", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".tpl.md") {
			newTpl, err := readTemplate(l.FS, path, tags)
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
