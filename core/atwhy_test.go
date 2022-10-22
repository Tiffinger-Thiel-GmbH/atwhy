package core_test

import (
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/finder"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/generator"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
)

func Test_core(t *testing.T) {
	// just full stack test
	wd, err := os.Getwd()
	assert.NoError(t, err)
	atwhy, err := core.New(generator.Markdown{}, filepath.Join(wd, ".."), "/", "templates", nil, map[string]finder.CommentConfig{
		".go": {
			LineComment: []string{"//"},
			BlockStart:  []string{"/*"},
			BlockEnd:    []string{"*/"},
		},
	})
	assert.NoError(t, err)

	templates, err := atwhy.Load()
	assert.NoError(t, err)
	assert.Greater(t, len(templates), 0)

	for _, tpl := range templates {
		fakeFs := afero.NewMemMapFs()
		file, err := fakeFs.Create(tpl.Path)
		assert.NoError(t, err)

		err = atwhy.Generate(tpl, file)
		assert.NoError(t, err)

		_, err = file.Seek(0, io.SeekStart)
		assert.NoError(t, err)

		content, err := afero.ReadAll(file)
		assert.NoError(t, err)

		// Test if there is actually something written to it.
		assert.Greater(t, len(content), 10)
	}
}
