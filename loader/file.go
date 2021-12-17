package loader

import (
	"io"
	"io/fs"
	"strings"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/tag"
	"github.com/aligator/checkpoint"
	"github.com/aligator/nogo"
	"github.com/spf13/afero"
)

type TagFinder interface {
	Find(filename string, reader io.Reader) (tags []tag.Raw, err error)
}

type File struct {
	FS             afero.Fs
	FileExtensions []string
}

// @WHY readme.ignore
// If you want to ignore files, just add a `.atwhyignore` to the root of your project.
// It follows the syntax of a `.gitignore` and you may also add `.atwhyignore` files to subfolders.

func (fl File) Load(finder TagFinder) ([]tag.Raw, error) {
	allTags := make([]tag.Raw, 0)

	n := nogo.New(nogo.DotGitRule)

	err := afero.Walk(fl.FS, ".", func(path string, info fs.FileInfo, err error) error {
		if ok, err := n.WalkFunc(afero.NewIOFS(fl.FS), ".atwhyignore", path, info.IsDir(), err); !ok {
			return checkpoint.From(err)
		}

		if info.IsDir() {
			return nil
		}

		if len(fl.FileExtensions) > 0 {
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
		}

		file, err := fl.FS.Open(path)
		if err != nil {
			return checkpoint.From(err)
		}
		defer file.Close()

		tags, err := finder.Find(path, file)
		if err != nil {
			return checkpoint.From(err)
		}
		allTags = append(allTags, tags...)
		return nil
	})
	if err != nil {
		return nil, checkpoint.From(err)
	}

	return allTags, nil
}
