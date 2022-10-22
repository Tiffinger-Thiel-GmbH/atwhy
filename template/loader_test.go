package template

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/afero"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/stretchr/testify/assert"
)

type fakeTag struct {
	name  string
	value string
}

func (f fakeTag) Type() tag.Type {
	panic("not implemented")
}

func (f fakeTag) String() string {
	return f.value
}

func (f fakeTag) Placeholder() string {
	return f.name
}

func Test_createTagMap(t *testing.T) {
	type args struct {
		tags []tag.Tag
	}
	tests := []struct {
		name string
		args args
		want mappedTags
	}{
		{
			name: "empty",
			args: args{
				tags: []tag.Tag{},
			},
			want: mappedTags{},
		},
		{
			name: "one",
			args: args{
				tags: []tag.Tag{
					fakeTag{"testName", "testValue"},
				},
			},
			want: mappedTags{
				"testName": fakeTag{"testName", "testValue"},
			},
		},
		{
			name: "some",
			args: args{
				tags: []tag.Tag{
					fakeTag{"testName", "testValue"},
					fakeTag{"testName2", "testValue2"},
					fakeTag{"testName3", "testValue3"},
				},
			},
			want: mappedTags{
				"testName":  fakeTag{"testName", "testValue"},
				"testName2": fakeTag{"testName2", "testValue2"},
				"testName3": fakeTag{"testName3", "testValue3"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := createTagMap(tt.args.tags); !reflect.DeepEqual(got, tt.want) {
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func testFS() afero.Fs {
	memFS := afero.NewMemMapFs()
	_ = memFS.MkdirAll("docu/subdocu", 0755)

	_ = afero.WriteFile(memFS, "README.tpl.md", []byte("# README"), 0777)
	_ = afero.WriteFile(memFS, "docu/Docu.tpl.md", []byte("# Docu"), 0777)
	_ = afero.WriteFile(memFS, "docu/subdocu/Docu2.tpl.md", []byte("# Docu2"), 0777)
	_ = afero.WriteFile(memFS, "LICENSE", []byte("GPL"), 0777)
	_ = afero.WriteFile(memFS, "OtherFile", []byte("Foo Bar"), 0777)
	_ = afero.WriteFile(memFS, "docu/OtherFile.md", []byte("* Foo Bar"), 0777)

	return memFS
}

func buildTestId(path string) string {
	id := md5.Sum([]byte(filepath.ToSlash(path)))
	return "page-" + hex.EncodeToString(id[:])
}

func TestLoader_Load(t *testing.T) {
	prefix := "/templates"

	type fields struct {
		FS                afero.Fs
		ProjectPathPrefix string
	}
	type args struct {
		tags []tag.Tag
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []Markdown
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "a bunch of files",
			fields: fields{
				FS:                testFS(),
				ProjectPathPrefix: prefix,
			},
			args: args{
				tags: []tag.Tag{},
			},
			want: []Markdown{
				{
					ID:                buildTestId("docu/Docu.tpl.md"),
					ProjectPathPrefix: prefix,
					Name:              "Docu",
					Path:              "docu",
					Header: Header{
						Meta: MetaData{Title: "Docu"},
					},
					tagMap: mappedTags{},
				},
				{
					ID:                buildTestId("docu/subdocu/Docu2.tpl.md"),
					ProjectPathPrefix: prefix,
					Name:              "Docu2",
					Path:              "docu/subdocu",
					Header: Header{
						Meta: MetaData{Title: "Docu2"},
					},
					tagMap: mappedTags{},
				},
				{
					ID:                buildTestId("README.tpl.md"),
					ProjectPathPrefix: prefix,
					Name:              "README",
					Path:              ".",
					Header: Header{
						Meta: MetaData{Title: "README"},
					},
					tagMap: mappedTags{},
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "no files",
			fields: fields{
				FS:                afero.NewMemMapFs(),
				ProjectPathPrefix: prefix,
			},
			args: args{
				tags: []tag.Tag{},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "broken fs",
			fields: fields{
				FS: afero.NewBasePathFs(testFS(), "NonExistent"), // This creates a broken fs.

				ProjectPathPrefix: prefix,
			},
			args: args{
				tags: []tag.Tag{},
			},
			want:    nil,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := Loader{
				FS:                tt.fields.FS,
				ProjectPathPrefix: tt.fields.ProjectPathPrefix,
			}
			got, err := l.Load(tt.args.tags)
			if !tt.wantErr(t, err, fmt.Sprintf("Load(%v)", tt.args.tags)) {
				return
			}

			for i := range tt.want {
				assert.NotNil(t, got[i].template)
				got[i].template = nil
			}
			assert.Equalf(t, tt.want, got, "Load(%v)", tt.args.tags)
		})
	}
}
