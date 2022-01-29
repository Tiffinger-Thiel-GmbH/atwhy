package template

import (
	"io/fs"
	"sort"
	"strings"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/spf13/afero"
)

type Loader struct {
	FS                afero.Fs
	ProjectPathPrefix string
}

type mappedTags = map[string]tag.Tag

func createTagMap(tags []tag.Tag) mappedTags {
	tagMap := make(map[string]tag.Tag)

	for _, t := range tags {
		tagMap[t.Placeholder()] = t
	}

	return tagMap
}

// Load templates from the Loader.FS.
func (l Loader) Load(tags []tag.Tag) ([]Markdown, error) {
	var res []Markdown

	mappedTags := createTagMap(tags)

	err := afero.Walk(l.FS, "", func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}

		if strings.HasSuffix(path, ".tpl.md") {
			newTpl, err := readTemplate(l.FS, l.ProjectPathPrefix, path, mappedTags)
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
		// Sort based on index-files, path depth
		if res[i].Header.Server.Index && !strings.HasPrefix(res[i].Path, res[j].Path) {
			return false
		}
		// and name
		return res[i].Header.Meta.Title < res[j].Header.Meta.Title
	})

	return res, nil
}
