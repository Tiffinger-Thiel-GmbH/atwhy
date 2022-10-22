package cmd

import (
	"testing"

	"github.com/Tiffinger-Thiel-GmbH/atwhy/finder"
	"github.com/stretchr/testify/assert"
)

func Test_generateCommentConfig(t *testing.T) {
	type args struct {
		comments []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]finder.CommentConfig
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "only line comments",
			args: args{
				comments: []string{
					"lol:#",
				},
			},
			want: map[string]finder.CommentConfig{
				".lol": {LineComment: []string{"#"}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "several line comments for one file type",
			args: args{
				comments: []string{
					"lol:#",
					"lol://",
				},
			},
			want: map[string]finder.CommentConfig{
				".lol": {LineComment: []string{"#", "//"}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "only block comments",
			args: args{
				comments: []string{
					"lol:,/*,*/",
				},
			},
			want: map[string]finder.CommentConfig{
				".lol": {BlockStart: []string{"/*"}, BlockEnd: []string{"*/"}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "missing block end",
			args: args{
				comments: []string{
					"lol:,/*",
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrInvalidCommentStringMissingBlock)
			},
		},
		{
			name: "missing block start",
			args: args{
				comments: []string{
					"lol:,,*/",
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrInvalidCommentStringMissingBlock)
			},
		},
		{
			name: "with , in the comments",
			args: args{
				comments: []string{
					`lol:\,,/*\,,\,*/`,
				},
			},
			want: map[string]finder.CommentConfig{
				".lol": {LineComment: []string{","}, BlockStart: []string{"/*,"}, BlockEnd: []string{",*/"}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "with : in the comments",
			args: args{
				comments: []string{
					`lol:::,:#,#:`,
				},
			},
			want: map[string]finder.CommentConfig{
				".lol": {LineComment: []string{"::"}, BlockStart: []string{":#"}, BlockEnd: []string{"#:"}},
			},
			wantErr: assert.NoError,
		},
		{
			name: "invalid base structure (missing :)",
			args: args{
				comments: []string{
					`la li lu`,
				},
			},
			want: nil,
			wantErr: func(tt assert.TestingT, err error, i ...interface{}) bool {
				return assert.ErrorIs(tt, err, ErrInvalidCommentString)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateCommentConfig(tt.args.comments)
			tt.wantErr(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}
