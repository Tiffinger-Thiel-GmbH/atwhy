package loader

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
	"io"
	"sort"
	"testing"
)

type fakeFile struct {
	Tags []tag.Raw `json:"Tags,omitempty"`
}

func (ff fakeFile) toJSON() []byte {
	result, err := json.Marshal(ff)
	if err != nil {
		panic(err)
	}

	return append([]byte("tags"), result...) // prefix with "tags" to indicate for the fake-finder to parse it.
}

var (
	testFileMainGo = fakeFile{
		Tags: []tag.Raw{
			{Type: tag.TypeDoc, Placeholder: "tag1", Filename: "main.go", Line: 1, Value: "value1"},
			{Type: tag.TypeDoc, Placeholder: "tag2", Filename: "main.go", Line: 10, Value: "value2"},
		},
	}
	testFileRunSh = fakeFile{
		Tags: []tag.Raw{
			{Type: tag.TypeDoc, Placeholder: "tag2", Filename: "run.sh", Line: 1, Value: "value2"},
			{Type: tag.TypeDoc, Placeholder: "tag3", Filename: "run.sh", Line: 10, Value: "value3"},
		},
	}
	testFileIgnored = fakeFile{
		Tags: []tag.Raw{
			{Type: tag.TypeDoc, Placeholder: "tag4", Filename: "ignored.go", Line: 1, Value: "value4"},
			{Type: tag.TypeDoc, Placeholder: "tag5", Filename: "ignored.go", Line: 10, Value: "value5"},
		},
	}
)

func testFs() afero.Fs {
	memFS := afero.NewMemMapFs()
	_ = afero.WriteFile(memFS, "main.go", testFileMainGo.toJSON(), 0777)
	_ = afero.WriteFile(memFS, "run.sh", testFileRunSh.toJSON(), 0777)
	_ = afero.WriteFile(memFS, "ignored.go", testFileIgnored.toJSON(), 0777)
	_ = afero.WriteFile(memFS, ".atwhyignore", []byte("/ignored.go"), 0777) // The github.com/aligator/nogo package is tested itself - no need to test it here deeply.
	return memFS
}

type fakeJSONFinder struct {
	returnError bool
}

func (f fakeJSONFinder) Find(_ string, reader io.Reader) ([]tag.Raw, error) {
	if f.returnError {
		return nil, errors.New("an error")
	}
	singleByte := make([]byte, 4)
	_, _ = reader.Read(singleByte)
	if !bytes.Equal(singleByte, []byte("tags")) {
		// Assume it is no json - just ignore
		return nil, nil
	}
	newFakeFile := fakeFile{}
	err := json.NewDecoder(reader).Decode(&newFakeFile)
	if err != nil {
		return nil, err
	}
	return newFakeFile.Tags, err
}

func TestFile_Load(t *testing.T) {
	type fields struct {
		FS             afero.Fs
		FileExtensions []string
	}
	type args struct {
		finder TagFinder
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []tag.Raw
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "find all but ignored",
			fields: fields{
				FS: testFs(),
			},
			args: args{
				finder: fakeJSONFinder{},
			},
			want:    append(testFileMainGo.Tags, testFileRunSh.Tags...),
			wantErr: assert.NoError,
		},
		{
			name: "find only go",
			fields: fields{
				FS:             testFs(),
				FileExtensions: []string{"go"},
			},
			args: args{
				finder: fakeJSONFinder{},
			},
			want:    testFileMainGo.Tags,
			wantErr: assert.NoError,
		},
		{
			name: "finder error gets handled",
			fields: fields{
				FS: testFs(),
			},
			args: args{
				finder: fakeJSONFinder{returnError: true},
			},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fl := File{
				FS:             tt.fields.FS,
				FileExtensions: tt.fields.FileExtensions,
			}
			got, err := fl.Load(tt.args.finder)
			if !tt.wantErr(t, err, "Load() error = %v", err) {
				return
			}

			// Order doesn't matter, but they have to exist all -> sort both the same way and then check.
			sort.SliceStable(got, func(i, j int) bool {
				return got[i].Placeholder > got[j].Placeholder
			})
			sort.SliceStable(tt.want, func(i, j int) bool {
				return tt.want[i].Placeholder > tt.want[j].Placeholder
			})

			assert.Equal(t, got, tt.want)
		})
	}
}
