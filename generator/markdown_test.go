package generator

import (
	"bytes"
	"errors"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/template"
	"github.com/stretchr/testify/assert"
	"testing"
)

type fakeErrorWriter struct{}

func (f fakeErrorWriter) Write(_ []byte) (n int, err error) {
	return 0, errors.New("an error")
}

func TestMarkdown_Ext(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "extension matches .md",
			want: ".md",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Markdown{}
			assert.Equal(t, tt.want, m.Ext())
		})
	}
}

func TestMarkdown_Generate(t *testing.T) {
	type args struct {
		template template.Markdown
	}
	tests := []struct {
		name       string
		args       args
		wantWriter string
		wantErr    assert.ErrorAssertionFunc
	}{
		{
			name: "executes the template and just returns the result (with an newline)",
			args: args{
				template: template.TestTemplate(template.Markdown{
					Header: template.Header{
						Meta: template.MetaData{Title: "supertemplate"},
					},
				}, `# Hello {{ .Meta.Title }}
* This
* Is
* Markdown`),
			},
			wantWriter: `# Hello supertemplate
* This
* Is
* Markdown
`,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := Markdown{}
			writer := &bytes.Buffer{}
			err := m.Generate(tt.args.template, writer)
			if !tt.wantErr(t, err, "Generate() error = %v", err) {
				return
			}
			assert.Equal(t, tt.wantWriter, writer.String())
		})
	}

	t.Run("nil buffer", func(t *testing.T) {
		m := Markdown{}
		writer := fakeErrorWriter{}
		err := m.Generate(template.TestTemplate(template.Markdown{}, `# Hello {{ .Meta.Title }}`), writer)
		assert.Error(t, err, "should throw error")
	})
}
