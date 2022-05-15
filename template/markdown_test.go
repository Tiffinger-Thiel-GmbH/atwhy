package template

import (
	"fmt"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"testing"
)

func testFileFS(name string, data []byte) afero.Fs {
	memFS := afero.NewMemMapFs()
	_ = afero.WriteFile(memFS, name, data, 0777)
	return memFS
}

func Test_readTemplate(t *testing.T) {
	prefix := "/template/"
	type args struct {
		sysfs             afero.Fs
		projectPathPrefix string
		path              string
		tags              mappedTags
	}
	tests := []struct {
		name            string
		args            args
		want            Markdown
		wantNilTemplate bool
		wantErr         assert.ErrorAssertionFunc
	}{
		{
			name: "simple file",
			args: args{
				sysfs:             testFileFS("README.tpl.md", []byte("# Hello")),
				projectPathPrefix: prefix,
				path:              "README.tpl.md",
			},
			want: Markdown{
				ID:                buildTestId("README.tpl.md"),
				ProjectPathPrefix: prefix,
				Name:              "README",
				Path:              ".",
				Value:             "# Hello",
				Header: Header{
					Meta: MetaData{Title: "README"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "with header",
			args: args{
				sysfs: testFileFS("README.tpl.md", []byte(`---
meta:
  title: Test Readme
server:
  index: true
---
# Hello World!`)),
				projectPathPrefix: prefix,
				path:              "README.tpl.md",
			},
			want: Markdown{
				ID:                buildTestId("README.tpl.md"),
				ProjectPathPrefix: prefix,
				Name:              "README",
				Path:              ".",
				Value:             "# Hello World!",
				Header: Header{
					Meta:   MetaData{Title: "Test Readme"},
					Server: ServerData{Index: true},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "empty file",
			args: args{
				sysfs:             testFileFS("README.tpl.md", nil),
				projectPathPrefix: prefix,
				path:              "README.tpl.md",
			},
			want: Markdown{
				ID:                buildTestId("README.tpl.md"),
				ProjectPathPrefix: prefix,
				Name:              "README",
				Path:              ".",
				Header: Header{
					Meta: MetaData{Title: "README"},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "non existent file",
			args: args{
				sysfs:             testFileFS("README.tpl.md", nil),
				projectPathPrefix: prefix,
				path:              "README2.tpl.md",
			},
			want:            Markdown{},
			wantNilTemplate: true,
			wantErr:         assert.Error,
		},
		{
			name: "invalid yaml header",
			args: args{
				sysfs: testFileFS("README.tpl.md", []byte(`---
meta:
- no
- valid
- yaml
- foo
- bar
---
# Hello World!`)),
				projectPathPrefix: prefix,
				path:              "README.tpl.md",
			},
			want:            Markdown{},
			wantNilTemplate: true,
			wantErr:         assert.Error,
		},
		{
			name: "invalid header - missing 2nd line",
			args: args{
				sysfs: testFileFS("README.tpl.md", []byte(`---
meta:
  title: Test Readme
server:
  index: true

# Hello World!`)),
				projectPathPrefix: prefix,
				path:              "README.tpl.md",
			},
			want:            Markdown{},
			wantNilTemplate: true,
			wantErr: func(t assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(t, err, ErrMissingBody, i)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readTemplate(tt.args.sysfs, tt.args.projectPathPrefix, tt.args.path, tt.args.tags)
			if !tt.wantErr(t, err, fmt.Sprintf("readTemplate(%v, %v, %v, %v)", tt.args.sysfs, tt.args.projectPathPrefix, tt.args.path, tt.args.tags)) {
				return
			}

			if !tt.wantNilTemplate {
				assert.NotNil(t, got.template, "template should not be nil")
				got.template = nil
			}

			assert.Equalf(t, tt.want, got, "readTemplate(%v, %v, %v, %v)", tt.args.sysfs, tt.args.projectPathPrefix, tt.args.path, tt.args.tags)
		})
	}
}
