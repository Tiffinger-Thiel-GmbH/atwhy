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
