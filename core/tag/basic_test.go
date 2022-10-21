package tag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_textFactory(t *testing.T) {
	type args struct {
		input      Raw
		isMarkdown bool
	}
	tests := []struct {
		name string
		args args
		want Basic
	}{
		{
			name: "just text",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file",
					Line:        5,
					Value:       "header\nvalue",
				},
				isMarkdown: false,
			},
			want: Basic{
				tagType:     TypeDoc,
				value:       "value",
				placeholder: "a_placeholder",
			},
		},
		{
			name: "cleans windows line endings",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file",
					Line:        5,
					Value:       "header\r\nvalue\r\nlol",
				},
				isMarkdown: false,
			},
			want: Basic{
				tagType:     TypeDoc,
				value:       "value\nlol",
				placeholder: "a_placeholder",
			},
		},
		{
			name: "one line only",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file",
					Line:        5,
					Value:       "header",
				},
				isMarkdown: false,
			},
			want: Basic{
				tagType:     TypeDoc,
				value:       "",
				placeholder: "a_placeholder",
			},
		},
		{
			name: "removes spaces and newlines at the end of the value",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file",
					Line:        5,
					Value:       "header\nfoo\nbar \n",
				},
				isMarkdown: false,
			},
			want: Basic{
				tagType:     TypeDoc,
				value:       "foo\nbar",
				placeholder: "a_placeholder",
			},
		},
		{
			name: "if it is markdown, it replaces all line endings by line endings with two space to force hard newlines",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file",
					Line:        5,
					Value:       "header\nfoo \nbar   \nbaz\nbum \n",
				},
				isMarkdown: true,
			},
			want: Basic{
				tagType:     TypeDoc,
				value:       "foo   \nbar     \nbaz  \nbum", // Note that for now its ok to just have more spaces. They don't hurt...
				placeholder: "a_placeholder",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := textFactory(tt.args.input, tt.args.isMarkdown)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestProjectLink(t *testing.T) {
	type args struct {
		input Raw
	}
	tests := []struct {
		name    string
		args    args
		want    Tag
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "a simple link",
			args: args{
				input: Raw{
					Type:        TypeLink,
					Placeholder: "a_placeholder",
					Filename:    "file.txt",
					Line:        5,
				},
			},
			want: Basic{
				tagType:     TypeLink,
				placeholder: "a_placeholder",
				value:       `[file.txt:5]({{ .Project "file.txt" }})`,
			},
			wantErr: assert.NoError,
		},
		{
			name: "not a TypeLink - should return nil, nil",
			args: args{
				input: Raw{
					Type:        TypeCode,
					Placeholder: "a_placeholder",
					Filename:    "file.txt",
					Line:        5,
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: `" should be escaped`,
			args: args{
				input: Raw{
					Type:        TypeLink,
					Placeholder: "a_placeholder",
					Filename:    `fi"le.txt`,
					Line:        5,
				},
			},
			want: Basic{
				tagType:     TypeLink,
				placeholder: "a_placeholder",
				value:       `[fi"le.txt:5]({{ .Project "fi\"le.txt" }})`,
			},
			wantErr: assert.NoError,
		},
		{
			name: `parentheses should be escaped`,
			args: args{
				input: Raw{
					Type:        TypeLink,
					Placeholder: "a_placeholder",
					Filename:    `fi(l)[e].txt`,
					Line:        5,
				},
			},
			want: Basic{
				tagType:     TypeLink,
				placeholder: "a_placeholder",
				value:       `[fi(l)\[e\].txt:5]({{ .Project "fi(l\)[e].txt" }})`,
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ProjectLink(tt.args.input)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestDoc(t *testing.T) {
	type args struct {
		input Raw
	}
	tests := []struct {
		name    string
		args    args
		want    Tag
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "simple text",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file.txt",
					Line:        5,
					Value:       "header\nsome\ntext\n",
				},
			},
			want: Basic{
				tagType:     TypeDoc,
				placeholder: "a_placeholder",
				value:       "some  \ntext",
			},
			wantErr: assert.NoError,
		},
		{
			name: "not a TypeDoc - should return nil, nil",
			args: args{
				input: Raw{
					Type:        TypeCode,
					Placeholder: "a_placeholder",
					Filename:    "file.txt",
					Line:        5,
					Value:       "header\nsome\ntext\n",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Doc(tt.args.input)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestCode(t *testing.T) {
	type args struct {
		input Raw
	}
	tests := []struct {
		name    string
		args    args
		want    Tag
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "simple text",
			args: args{
				input: Raw{
					Type:        TypeCode,
					Placeholder: "a_placeholder",
					Filename:    "file.txt",
					Line:        5,
					Value:       "header\nsome\ntext\n",
				},
			},
			want: Basic{
				tagType:     TypeCode,
				placeholder: "a_placeholder",
				value:       "```txt\nsome\ntext\n```\n",
			},
			wantErr: assert.NoError,
		},
		{
			name: "not a TypeCode - should return nil, nil",
			args: args{
				input: Raw{
					Type:        TypeDoc,
					Placeholder: "a_placeholder",
					Filename:    "file.txt",
					Line:        5,
					Value:       "header\nsome\ntext\n",
				},
			},
			want:    nil,
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Code(tt.args.input)

			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
