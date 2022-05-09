package cmd

import (
	"errors"
	"github.com/Tiffinger-Thiel-GmbH/atwhy/finder"
	"github.com/spf13/cobra"
	"path/filepath"
	"strings"
)

var defaultComments = []string{
	// @WHY CODE readme_comments_builtin
	"://,/*,*/",
	"sh:#",
	`py:#,""","""`,
	"cmd:REM",
	`cmd:::`,
	"vb:'",
	"html,xml:,<!--,-->",
	"lua:--,--[[,]]",
	"sql:--",
	"rb:#,=begin,=end",
	// @WHY CODE_END
}

var ErrInvalidCommentString = errors.New("comment configuration has to be like '{extList}:{lineComment},{blockStart},{blockEnd}' (see --help)")

// LoadCommonArgs loads everything which is common through the different modes.
func LoadCommonArgs(cmd *cobra.Command) (templatesFolder string, projectPath string, extensions []string, commentConfig map[string]finder.CommentConfig, err error) {
	templatesFolder, err = cmd.Flags().GetString("templates-folder")
	if err != nil {
		return "", "", nil, nil, err
	}

	projectPath, err = cmd.Flags().GetString("project")
	if err != nil {
		return "", "", nil, nil, err
	}

	// Make the path absolute.
	projectPath, err = filepath.Abs(projectPath)
	if err != nil {
		return "", "", nil, nil, err
	}

	extensions, err = cmd.Flags().GetStringSlice("ext")
	if err != nil {
		return "", "", nil, nil, err
	}

	comments, err := cmd.Flags().GetStringArray("comment")
	if err != nil {
		return "", "", nil, nil, err
	}

	commentConfig = make(map[string]finder.CommentConfig)
	if len(comments) == 0 {
		comments = append(comments, "DEFAULT")
	}

	// Replace DEFAULT with the default comments.
	extendedComments := make([]string, 0, len(comments))
	for i, comment := range comments {
		if comment == "DEFAULT" {
			preComments := comments[:i]
			extendedComments = append(preComments, defaultComments...)
		} else {
			extendedComments = append(extendedComments, comment)
		}
	}

	comments = extendedComments

	for _, comment := range comments {
		split := strings.SplitN(comment, ":", 2)
		if len(split) != 2 {
			return "", "", nil, nil, ErrInvalidCommentString
		}

		commaPlaceholder := string([]byte{1})

		cfgSplit := strings.Split(strings.ReplaceAll(split[1], `\,`, commaPlaceholder), ",")

		commentExtensions := strings.Split(split[0], ",")
		for _, ext := range commentExtensions {
			ext = "." + ext

			if _, ok := commentConfig[ext]; !ok {
				commentConfig[ext] = finder.CommentConfig{}
			}

			cfg := commentConfig[ext]

			if len(cfgSplit) >= 1 && cfgSplit[0] != "" {
				cfg.LineComment = append(commentConfig[ext].LineComment, cfgSplit[0])
			}

			if len(cfgSplit) >= 3 && cfgSplit[1] != "" && cfgSplit[2] != "" {
				cfg.BlockStart = append(commentConfig[ext].BlockStart, cfgSplit[1])
				cfg.BlockEnd = append(commentConfig[ext].BlockEnd, cfgSplit[2])
			}

			commentConfig[ext] = cfg
		}
	}

	return templatesFolder, projectPath, extensions, commentConfig, nil
}
