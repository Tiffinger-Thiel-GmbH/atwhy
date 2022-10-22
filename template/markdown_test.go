package template

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/spf13/afero"
	"github.com/stretchr/testify/assert"
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

func Test_data_Project(t *testing.T) {
	type fields struct {
		Tag              map[string]tag.Tag
		Meta             MetaData
		Now              string
		projectPrefix    string
		isPostprocessing bool
	}
	type args struct {
		file string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "joins prefix with filename",
			fields: fields{
				projectPrefix: "prefix",
			},
			args: args{
				file: "filename.md",
			},
			want: "prefix/filename.md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := data{
				Tag:              tt.fields.Tag,
				Meta:             tt.fields.Meta,
				Now:              tt.fields.Now,
				projectPrefix:    tt.fields.projectPrefix,
				isPostprocessing: tt.fields.isPostprocessing,
			}
			assert.Equalf(t, tt.want, d.Project(tt.args.file), "Project(%v)", tt.args.file)
		})
	}
}

func Test_data_Escape(t *testing.T) {
	type fields struct {
		Tag              map[string]tag.Tag
		Meta             MetaData
		Now              string
		projectPrefix    string
		isPostprocessing bool
	}
	type args struct {
		value string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "no postprocessing",
			fields: fields{
				isPostprocessing: false,
			},
			args: args{
				value: "{{ .SomeValue }}",
			},
			want: `{{"{{ .SomeValue }}"}}`,
		},
		{
			name: "with \\\"",
			fields: fields{
				isPostprocessing: false,
			},
			args: args{
				value: `{{ \".SomeValue\" }}`,
			},
			want: `{{"{{ \\\".SomeValue\\\" }}"}}`,
		},
		{
			name: "with not escaped \"",
			fields: fields{
				isPostprocessing: false,
			},
			args: args{
				value: `{{ ".SomeValue" }}`,
			},
			want: `{{"{{ \".SomeValue\" }}"}}`,
		},
		{
			name: "with postprocessing - no actual escaping needed",
			fields: fields{
				isPostprocessing: true,
			},
			args: args{
				value: "{{ .SomeValue }}",
			},
			want: "{{ .SomeValue }}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := data{
				Tag:              tt.fields.Tag,
				Meta:             tt.fields.Meta,
				Now:              tt.fields.Now,
				projectPrefix:    tt.fields.projectPrefix,
				isPostprocessing: tt.fields.isPostprocessing,
			}
			assert.Equalf(t, tt.want, d.Escape(tt.args.value), "Escape(%v)", tt.args.value)
		})
	}
}

func TestMarkdown_Execute(t1 *testing.T) {
	type fields struct {
		ID                string
		ProjectPathPrefix string
		Name              string
		Path              string
		Value             string
		Header            Header
		template          *template.Template
		tagMap            map[string]tag.Tag
	}
	tests := []struct {
		name       string
		fields     fields
		wantWriter string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "Render .Meta.Title",
			fields: fields{
				Header: Header{
					Meta: MetaData{Title: "Readme"},
				},
				template: template.Must(template.New("").Parse("Title: {{ .Meta.Title }}")),
				tagMap:   make(map[string]tag.Tag),
			},
			wantWriter: "Title: Readme",
			wantErr:    assert.NoError,
		},
		{
			name: "Render .Tag.something",
			fields: fields{
				template: template.Must(template.New("").Parse("Tag: {{ .Tag.something }}")),
				tagMap: map[string]tag.Tag{
					"something": fakeTag{
						name:  "something",
						value: "different",
					},
				},
			},
			wantWriter: "Tag: different",
			wantErr:    assert.NoError,
		},
		{
			name: "Render .Tag with another, nested placeholder",
			fields: fields{
				Header: Header{
					Meta: MetaData{Title: "Readme"},
				},
				template: template.Must(template.New("").Parse("Title: {{ .Tag.something }}")),
				tagMap: map[string]tag.Tag{
					"something": fakeTag{
						name:  "something",
						value: "\"{{ .Meta.Title }}\"",
					},
				},
			},
			wantWriter: "Title: \"Readme\"",
			wantErr:    assert.NoError,
		},
		{
			name: "Render .Escape",
			fields: fields{
				template: template.Must(template.New("").Parse(`Escape: {{ .Escape "{{ \"huhu\" }}" }}`)),
				tagMap:   make(map[string]tag.Tag),
			},
			wantWriter: `Escape: {{ "huhu" }}`,
			wantErr:    assert.NoError,
		},
		{
			name: "Render .Escape in tag",
			fields: fields{
				template: template.Must(template.New("").Parse(`Escape: {{ .Tag.something }}`)),
				tagMap: map[string]tag.Tag{
					"something": fakeTag{
						name:  "something",
						value: `something escaped: {{ .Escape "{{ \"huhu\" }}" }}`,
					},
				},
			},
			wantWriter: `Escape: something escaped: {{ "huhu" }}`,
			wantErr:    assert.NoError,
		},
		{
			name: "non existent tag",
			fields: fields{
				template: template.Must(template.New("").Parse(`NOO: {{ .Tag.none }}`)),
				tagMap: map[string]tag.Tag{
					"something": fakeTag{
						name:  "something",
						value: `different`,
					},
				},
			},
			wantWriter: `NOO: <no value>`,
			wantErr:    assert.NoError,
		},
		{
			name: "tag with invalid template",
			fields: fields{
				template: template.Must(template.New("").Parse(`NOO: {{ .Tag.something }}`)),
				tagMap: map[string]tag.Tag{
					"something": fakeTag{
						name:  "something",
						value: `{{ .broken.`,
					},
				},
			},
			wantErr: assert.Error,
		},
		{
			name: "tag with invalid template - escaped",
			fields: fields{
				template: template.Must(template.New("").Parse(`NOO: {{ .Tag.something }}`)),
				tagMap: map[string]tag.Tag{
					"something": fakeTag{
						name:  "something",
						value: `{{ .Escape "{{ .broken." }}`,
					},
				},
			},
			wantWriter: "NOO: {{ .broken.",
			wantErr:    assert.NoError,
		},
	}
	for _, tt := range tests {
		t1.Run(tt.name, func(t1 *testing.T) {
			t := Markdown{
				ID:                tt.fields.ID,
				ProjectPathPrefix: tt.fields.ProjectPathPrefix,
				Name:              tt.fields.Name,
				Path:              tt.fields.Path,
				Header:            tt.fields.Header,
				template:          tt.fields.template,
				tagMap:            tt.fields.tagMap,
			}
			writer := &bytes.Buffer{}
			err := t.Execute(writer)
			if !tt.wantErr(t1, err, fmt.Sprintf("Execute(%v)", writer)) {
				return
			}
			assert.Equalf(t1, tt.wantWriter, writer.String(), "Execute(%v)", writer)
		})
	}
}

func mustTag(t tag.Tag, err error) tag.Tag {
	if err != nil {
		panic(err)
	}
	return t
}

func Test_data_Group(t *testing.T) {
	type fields struct {
		Tag              map[string]tag.Tag
		Meta             MetaData
		Now              string
		projectPrefix    string
		isPostprocessing bool
	}
	type args struct {
		prefix string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   string
	}{
		{
			name: "group some tags",
			fields: fields{
				Tag: map[string]tag.Tag{
					"test-0": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "test-0",
						Value:       "header\na test0",
					})),
					"noo-1": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "noo-1",
						Value:       "header\na test-1",
					})),
					"test-1": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "test-1",
						Value:       "header\na test1",
					})),
				},
			},
			args: args{"test"},
			want: `a test0  
a test1  
`,
		},
		{
			name: "the tags get sorted alphanumerical",
			fields: fields{
				Tag: map[string]tag.Tag{
					"testa-1": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "testa-1",
						Value:       "header\na-1",
					})),
					"testa-100": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "testa-100",
						Value:       "header\na-100",
					})),
					"testa-111": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "testa-111",
						Value:       "header\na-111",
					})),
					"testa-10": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "testa-10",
						Value:       "header\na-10",
					})),
					"testb-0": mustTag(tag.Doc(tag.Raw{
						Type:        tag.TypeDoc,
						Placeholder: "testb-0",
						Value:       "header\nb-0",
					})),
				},
			},
			args: args{"test"},
			want: `a-1  
a-10  
a-100  
a-111  
b-0  
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := data{
				Tag:              tt.fields.Tag,
				Meta:             tt.fields.Meta,
				Now:              tt.fields.Now,
				projectPrefix:    tt.fields.projectPrefix,
				isPostprocessing: tt.fields.isPostprocessing,
			}
			assert.Equal(t, tt.want, d.Group(tt.args.prefix))
		})
	}
}
