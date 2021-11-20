package loader

import (
	"fmt"
	ignore "github.com/sabhiram/go-gitignore"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/afero"
	"gitlab.com/tiffinger-thiel/crazydoc/tag"
)

type TagFinder interface {
	Find(filename string, reader io.Reader) (tags []tag.Raw, err error)
}

type File struct {
	FS             afero.Fs
	FileExtensions []string
}

// @README 20
// Ignore
// If you want to ignore files, just add a `.crazydocignore`
// to the root of your project. It follows the syntax of a `.gitignore`.
//
// __NOTE:__ Currently only the file in the root directory is checked, no
// deeper nested file, like with `.gitignore`.

func loadIgnores(path string) (*ignore.GitIgnore, error) {
	// Ignore .git always.
	rules := ignore.CompileIgnoreLines(".git")

	ignoreFile := filepath.Join(path, ".crazydocignore")
	if _, err := os.Stat(ignoreFile); err != nil {
		if os.IsNotExist(err) {
			return rules, nil
		}
		return rules, err
	}

	fmt.Println(ignoreFile)
	ignores, err := ignore.CompileIgnoreFileAndLines(ignoreFile)
	return ignores, err
}

func (fl File) Load(dir string, finder TagFinder) ([]tag.Raw, error) {
	filesystem := fl.FS
	allTags := make([]tag.Raw, 0)

	ignores, err := loadIgnores(dir)
	if err != nil {
		return nil, err
	}

	err = afero.Walk(fl.FS, dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if ignores.MatchesPath(path) {
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
