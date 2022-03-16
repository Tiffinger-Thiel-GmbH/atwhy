package finder

import (
	"github.com/Tiffinger-Thiel-GmbH/atwhy/core/tag"
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

var testCommentConfig = map[string]CommentConfig{
	".go": {
		LineComment: []string{"//"},
		BlockStart:  []string{"/*"},
		BlockEnd:    []string{"*/"},
	},
}

func TestFinder_Find(t *testing.T) {
	type fields struct {
		CommentConfig map[string]CommentConfig
	}
	type args struct {
		filename string
		reader   io.Reader
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []tag.Raw
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "normal @WHY tags",
			fields: fields{
				CommentConfig: testCommentConfig,
			},
			args: args{
				filename: "file.go",
				reader: strings.NewReader(`This is some fil
// @WHY my_tag_name
// Some text
//
// Another text block

// Separate comment block
// ...

/* 
@WHY my_other_tag_name
A Block comment.
...


Block end
*/
No part of comment`),
			},
			want: []tag.Raw{
				{
					Type:        tag.TypeDoc,
					Placeholder: "my_tag_name",
					Filename:    "file.go",
					Line:        1,
					Value: `@WHY my_tag_name
Some text

Another text block
`,
				},
				{
					Type:        tag.TypeDoc,
					Placeholder: "my_other_tag_name",
					Filename:    "file.go",
					Line:        10,
					Value: `@WHY my_other_tag_name
A Block comment.
...


Block end
`,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "filetype not in commentConfig",
			fields: fields{
				CommentConfig: testCommentConfig,
			},
			args: args{
				filename: "file.zz",
				reader: strings.NewReader(`This is some fil
// @WHY my_tag_name
// Some text
//
// Another text block

// Separate comment block
// ...

/* 
@WHY my_other_tag_name
A Block comment.
...


Block end
*/
No part of comment`),
			},
			want:    nil,
			wantErr: assert.NoError,
		},
		{
			name: "code tag",
			fields: fields{
				CommentConfig: testCommentConfig,
			},
			args: args{
				filename: "file.go",
				reader: strings.NewReader(`This is some fil
// @WHY CODE my_tag_name
// Some text
//
// Another text block
COODDDEEE
adsfj
asd

// Separate comment block
// ...
// @WHY CODE_END
NOT included

`),
			},
			want: []tag.Raw{
				{
					Type:        tag.TypeCode,
					Placeholder: "my_tag_name",
					Filename:    "file.go",
					Line:        1,
					Value: `@WHY CODE my_tag_name
// Some text
//
// Another text block
COODDDEEE
adsfj
asd

// Separate comment block
// ...
`,
				},
			},
			wantErr: assert.NoError,
		},
		{
			name: "link tag",
			fields: fields{
				CommentConfig: testCommentConfig,
			},
			args: args{
				filename: "file.go",
				reader: strings.NewReader(`This is some fil
// @WHY LINK my_link_tag
// Some text

// @WHY LINK another_link_tag
`),
			},
			want: []tag.Raw{
				{
					Type:        tag.TypeLink,
					Placeholder: "my_link_tag",
					Filename:    "file.go",
					Line:        1,
					Value: `@WHY LINK my_link_tag
`,
				},
				{
					Type:        tag.TypeLink,
					Placeholder: "another_link_tag",
					Filename:    "file.go",
					Line:        4,
					Value: `@WHY LINK another_link_tag
`,
				},
			},
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Finder{
				CommentConfig: tt.fields.CommentConfig,
			}
			got, err := f.Find(tt.args.filename, tt.args.reader)
			if !tt.wantErr(t, err) {
				return
			}

			assert.EqualValues(t, tt.want, got)
		})
	}
}
